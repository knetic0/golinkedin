package golinkedin

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Config struct {
	// The response type, set to 'code' by default, which is the standard for authorization flows.
	ResponseType string

	// Your LinkedIn app's ClientID, which is obtained after creating a LinkedIn app.
	ClientID string

	// Your LinkedIn app's ClientSecret, obtained after creating the LinkedIn app.
	ClientSecret string

	// The callback URL you define when setting up your LinkedIn app. It will be used for redirecting users.
	RedirectURI string

	// The requested permissions (scopes). Set to 'r_liteprofile', 'r_emailaddress', 'w_member_social', and 'w_share' by default.
	Scopes []string
}

type Linkedin struct {
	// The configuration for the LinkedIn OAuth client, using the `Config` struct.
	config Config
}

var (
	// Base URL for LinkedIn's OAuth authentication endpoint.
	AuthenticationBaseURL = "https://www.linkedin.com/oauth/v2/authorization?"

	// Base URL for requesting the access token after obtaining the authorization code.
	AccessTokenBaseURL = "https://www.linkedin.com/uas/oauth2/accessToken?grant_type=authorization_code"
)

/*
New creates a new LinkedIn API client with the given configuration.
This requires the Client ID, Redirect URL, and Client Secret, which can be found
on the LinkedIn Developer portal. Optionally, permissions (scopes) can be provided.
If no config is provided, an error is returned.
*/
func New(config ...Config) (*Linkedin, error) {
	linkedin := &Linkedin{
		config: Config{},
	}

	if len(config) == 0 {
		return nil, fmt.Errorf("%s", "please provide a config!")
	}

	linkedin.config = config[0]

	return linkedin, nil
}

/*
setAuthURL builds the LinkedIn authorization URL for the client.
It generates a random state for security and joins the required scopes.
The resulting URL is used to redirect the user for authentication.
*/
func (config *Config) setAuthURL() string {
	if len(config.Scopes) == 0 {
		config.Scopes = []string{"r_liteprofile", "r_emailaddress", "w_member_social", "w_share"}
	}

	state := stateGenerator()
	scopes := strings.Join(config.Scopes, ",")

	return AuthenticationBaseURL +
		"response_type=" + config.ResponseType +
		"&client_id=" + config.ClientID +
		"&redirect_uri=" + config.RedirectURI +
		"&state=" + state +
		"&scope=" + scopes +
		"&client_secret=" + config.ClientSecret
}

/*
setAccessTokenURL generates the access token request URL, which is used to exchange
the authorization code for an access token. It includes the redirect URI, client ID, and client secret.
*/
func (cfg *Config) setAccessTokenURL(code string) string {
	return AccessTokenBaseURL + code +
		"&redirect_uri=" + cfg.RedirectURI +
		"&client_id=" + cfg.ClientID +
		"&client_secret=" + cfg.ClientSecret
}

/*
GetAuthenticationUrl returns the LinkedIn authorization URL, allowing the client to
authenticate users and receive an authorization code.
*/
func (ln *Linkedin) GetAuthenticationUrl() string {
	return ln.config.setAuthURL()
}

/*
RetrieveAccessToken exchanges the authorization code for an access token by making an
HTTP request to LinkedIn's OAuth API. The access token is required for accessing LinkedIn's API.
*/
func (ln *Linkedin) RetrieveAccessToken(code string) (string, error) {
	type response struct {
		AccessToken string `json:"access_token"`
	}

	client := http.Client{}

	accessTokenURL := ln.config.setAccessTokenURL(code)

	req, _ := http.NewRequest(http.MethodGet, accessTokenURL, nil)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request error: %s", err.Error())
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body error: %s", err.Error())
	}

	var r response
	_ = json.Unmarshal(data, &r)

	return r.AccessToken, nil
}

/*
stateGenerator generates a random state string for OAuth2 requests.
This state is used to prevent CSRF attacks during the authentication process.
*/
func stateGenerator() string {
	b := make([]byte, 20)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
