package golinkedin

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Authentication interface {
	Authorize() string
	RedirectURL(*fiber.Ctx) error
	RetrieveAccessToken(*fiber.Ctx) error
	Profile(*fiber.Ctx) error
	SharePost(*fiber.Ctx) error
	ShareJOBPosting(*fiber.Ctx) error
}

var (
	StateToken = session.New(session.Config{
		KeyLookup: "cookie:state",
	})
	AccessToken = session.New(session.Config{
		KeyLookup: "cookie:access_token",
	})
	AuthURL = "https://www.linkedin.com/oauth/v2"
)

/*
First, we call this function to define some necessary part.
Send credentials as type of map[string]string.
  - ClientID `client_id`
  - RedirectURI `redirect_uri`
  - ClientSecret `client_secret`
*/
func New(clientID, redirectURI, clientSecret string) *API {
	api := API{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		ResponseType: "code",
		Scope:        "r_liteprofile,r_emailaddress,w_member_social",
		State:        tokenGenerator(),
	}

	api.AuthURL = api.Authorize()

	return &api
}

/*
Create Response URL. This URL helping us about take our AccessToken.
*/
func (s *API) Authorize() string {
	response_url := fmt.Sprintf("%s/authorization?response_type=%s&client_id=%s&redirect_uri=%s&state=%s&scope=%s&client_secret=%s",
		AuthURL, s.ResponseType, s.ClientID, s.RedirectURI, s.State, s.Scope, s.ClientSecret,
	)

	return response_url
}

/*
Call this function after InitAuth.
This function will redirect you to AuthURL, and After then run Retrieve function as automatically.
Must create route for this function.
*/
func (s *API) RedirectURL(c *fiber.Ctx) error {
	sess, err := StateToken.Get(c)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "state cookie is not defined"})
	}

	sess.Set("state", s.State)
	if err = sess.Save(); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "state cookie is not saved"})
	}

	return c.Redirect(s.AuthURL, http.StatusFound)
}

/*
Create route for this function.
You define Callback url as this route on your Linkedin app panel.
*/
func (s *API) RetrieveAccessToken(c *fiber.Ctx) error {
	queryState := c.Query("state")

	sess, err := StateToken.Get(c)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "state token is not fetching"})
	}

	if queryState != sess.Get("state") {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "authorization error, states are not same"})
	}

	client := http.Client{}

	queryToken := c.Query("code")

	accessTokenURL := "https://www.linkedin.com/uas/oauth2/accessToken?grant_type=authorization_code&code=" + queryToken + "&redirect_uri=" +
		s.RedirectURI + "&client_id=" + s.ClientID + "&client_secret=" + s.ClientSecret

	req, err := http.NewRequest(http.MethodGet, accessTokenURL, nil)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "get access token error"})
	}

	resp, err := client.Do(req)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "client do error"})
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "body read error"})
	}
	resp.Body.Close()

	var responseBody map[string]interface{}
	err = json.Unmarshal(data, &responseBody)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "unmarshal error"})
	}

	if _, err := responseBody["error"]; err {
		c.Status(400)
		return c.JSON(fiber.Map{"error": responseBody["error"].(string)})
	}

	s.AccessToken = responseBody["access_token"].(string)

	accessSession, err := AccessToken.Get(c)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "access token cookie error."})
	}

	accessSession.Set("access_token", s.AccessToken)
	if err = accessSession.Save(); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "session save error on access_token."})
	}

	c.Status(200)
	return c.JSON(fiber.Map{"message": "retrieve access token function worked as successfully"})
}

/*
tokenGenerator() function creating TOKEN for state query.
*/
func tokenGenerator() string {
	b := make([]byte, 20)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
