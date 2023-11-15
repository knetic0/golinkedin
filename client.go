package golinkedin

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Authentication interface {
	Authorize() string
	RedirectToURL(*fiber.Ctx) error
	RetrieveAccessToken(*fiber.Ctx) error
	Profile(*fiber.Ctx) error
	SharePost(*fiber.Ctx) error
	ShareJOBPosting(*fiber.Ctx) error
}

type Linkedin struct {
	// Redirect to linkedin Auth screen with this url.
	// Must be contain ClientID, RedirectURI, ClientSecret, Scopes and State.
	AuthURL string `json:"authURL"`

	// Set 'code' by default.
	ResponseType string `json:"responseType"`

	// Your linkedin app's ClientID.
	// You can take your ClientID after create linkedin app.
	ClientID string `json:"client_id"`

	// Your linkedin app's ClientSecret.
	// You can take your ClientSecret after create linkedin app.
	ClientSecret string `json:"client_secret"`

	// Your callback url.
	// You can define this while you creating your linkedin app.
	RedirectURI string `json:"redirect_uri"`

	// Your State token, creating random.
	State string `json:"state"`

	// Our permissions.
	// Set 'r_liteprofile,r_emailaddress,w_member_social,w_share' by defualt.
	Scope string `json:"scope"`

	// You can take AccessToken when you redirect to callback url.
	// This defined in your redirect URL as 'code'.
	AccessToken string `json:"access_token"`

	// ProfileInformation inherit on API Struct.
	ProfileInformation ProfileInformation `json:"profile_information"`
}

var (
	AuthURL = "https://www.linkedin.com/oauth/v2/authorization?"
)

/*
New function create a new linkedin API struct.
Take arguments, Client ID, Redirect URL and
Client Secret. Thats arguments can find on
linkedin portal. Also, send Scopes argument for
give permission your api.
*/
func New(clientId, redirectUrl, clientSecret string, scopes []string) (*Linkedin, error) {
	if len(scopes) == 0 {
		return nil, fmt.Errorf("scopes must take valid value")
	}

	for _, scp := range scopes {
		switch scp {
		case "r_liteprofile", "r_emailaddress", "w_member_social":
			continue
		default:
			return nil, fmt.Errorf("invalid scope")
		}
	}

	api := Linkedin{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURI:  redirectUrl,
		ResponseType: "code",
		Scope:        strings.Join(scopes, ","),
		State:        tokenGenerator(),
	}

	api.AuthURL = fmt.Sprintf("%sresponse_type=%s&client_id=%s&redirect_uri=%s&state=%s&scope=%s&client_secret=%s", AuthURL, api.ResponseType, api.ClientID, api.RedirectURI, api.State, api.Scope, api.ClientSecret)

	return &api, nil
}

/*
Redirect URL will be redirect you to AuthURL.
*/
func (ln *Linkedin) RedirectToURL(c *fiber.Ctx) error {
	return c.Redirect(ln.AuthURL, http.StatusFound)
}

/*
After redirecting to auth url, you must to callback here.
When you set something on linkedin portal you must to be
determine call back url, and this url point to RetrieveAccessToken's
endpoint.
*/
func (ln *Linkedin) RetrieveAccessToken(c *fiber.Ctx) error {
	client := http.Client{}

	queryToken := c.Query("code")

	accessTokenURL := "https://www.linkedin.com/uas/oauth2/accessToken?grant_type=authorization_code&code=" + queryToken + "&redirect_uri=" + ln.RedirectURI + "&client_id=" + ln.ClientID + "&client_secret=" + ln.ClientSecret

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

	ln.AccessToken = responseBody["access_token"].(string)

	c.Cookie(&fiber.Cookie{
		Name:     "linkedin_token",
		Value:    ln.AccessToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
	})

	return c.Status(200).JSON(fiber.Map{"access_token": ln.AccessToken})
}

/*
tokenGenerator() function creating TOKEN for state query.
*/
func tokenGenerator() string {
	b := make([]byte, 20)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
