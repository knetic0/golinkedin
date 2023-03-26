package golinkedin

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
	"testing"
)

func Routes() {
	app := fiber.New()

	app.Get("/", RedirectURL)
	app.Get("/callback", RetrieveAccessToken)
	app.Get("/profile", Profile)
	app.Get("/share", SharePost)
	app.Get("/sharejob", ShareJOBPosting)

	log.Fatal(app.Listen(":8000"))
}

func TestInitAuth(t *testing.T) {
	var credentials = map[string]string{
		"client_id":     "client_id",
		"redirect_uri":  "redirect_uri",
		"client_secret": "client_secret",
	}

	InitAuth(credentials)
	fmt.Println(Config.AuthURL)

	rUrl := Config.AuthURL

	if strings.Contains(rUrl, credentials["client_id"]) && strings.Contains(rUrl, credentials["redirect_uri"]) {
		t.Log("Success, URL is working...")
	}
}

func TestRedirectURL(t *testing.T) {
	var credentials = map[string]string{
		"client_id":     "77ljob4304lxf4",
		"redirect_uri":  "http://localhost:8000/callback",
		"client_secret": "0liGSEzhek8hjEdi",
	}

	InitAuth(credentials)

	Routes()
}
