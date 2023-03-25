package golinkedin

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"io/ioutil"
	"net/http"
)

type API struct {
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
}

var (
	Config API
	StateToken = session.New(session.Config{
		KeyLookup: "cookie:state",
	})
	AccessToken = session.New(session.Config{
		KeyLookup: "cookie:access_token",
	})
)

/*
	tokenGenerator() function creating TOKEN for state query.
*/
func tokenGenerator() string {
	b := make([]byte, 20)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

/*
	First, we call this function to define some necessary part.
	Send credentials as type of map[string]string.
	 	- ClientID `client_id`
		- RedirectURI `redirect_uri`
		- ClientSecret `client_secret`
 */
func InitAuth(credentials map[string]string) {
	Config.ClientID = credentials["client_id"]
	Config.RedirectURI = credentials["redirect_uri"]
	Config.ClientSecret = credentials["client_secret"]
	auth_url := "https://www.linkedin.com/oauth/v2"

	response_url := Authorize(auth_url)

	Config.AuthURL = response_url
}

/*
	Create Response URL. This URL helping us about take our AccessToken.
*/
func Authorize(auth_url string) string {
	csrf_token := tokenGenerator()

	Config.ResponseType = "code"
	Config.Scope = "r_liteprofile,r_emailaddress,w_member_social"
	Config.State = csrf_token

	response_url := fmt.Sprintf("%s/authorization?response_type=%s&client_id=%s&redirect_uri=%s&state=%s&scope=%s&client_secret=%s",
			auth_url, Config.ResponseType, Config.ClientID, Config.RedirectURI, Config.State, Config.Scope, Config.ClientSecret,
		)

	return response_url
}

/*
	Call this function after InitAuth.
	This function will redirect you to AuthURL, and After then run Retrieve function as automatically.
	Must create route for this function.
 */
func RedirectURL(c *fiber.Ctx) error {
	sess, err := StateToken.Get(c)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "state cookie is not defined"})
	}

	sess.Set("state", Config.State)
	if err = sess.Save(); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "state cookie is not saved"})
	}

	return c.Redirect(Config.AuthURL, http.StatusFound)
}

/*
	Create route for this function.
	You define Callback url as this route on your Linkedin app panel.
 */
func RetrieveAccessToken(c *fiber.Ctx) error {
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

	queryToken := c.Query("code")

	resp, err := http.Get("https://www.linkedin.com/uas/oauth2/accessToken?grant_type=authorization_code&code=" + queryToken + "&redirect_uri=" +
		Config.RedirectURI + "&client_id=" + Config.ClientID + "&client_secret=" + Config.ClientSecret)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "get access token error"})
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

	Config.AccessToken = responseBody["access_token"].(string)

	accessSession, err := AccessToken.Get(c)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error":"access token cookie error."})
	}

	accessSession.Set("access_token", Config.AccessToken)
	if err = accessSession.Save(); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "session save error on access_token."})
	}

	c.Status(200)
	return c.JSON(fiber.Map{"message": "retrieve access token function worked as successfully"})
}