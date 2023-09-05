package golinkedin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ProfileInformation struct {
	// ProfileID represents the ID every linkedin profile has.
	Id string `json:"id"`

	// User's FirstName
	FirstName string `json:"localizedFirstName"`

	// User's LastName
	LastName string `json:"localizedLastName"`
}

/*
After than Callback, you take Profile information with this function.
Create route for this.
*/
func (ln *Linkedin) Profile(c *fiber.Ctx) error {
	profile_url := "https://api.linkedin.com/v2/me"

	token := c.Cookies("linkedin_token")

	authorization := fmt.Sprintf("Bearer %s", token)

	client := http.Client{}

	req, err := http.NewRequest("GET", profile_url, nil)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "new request created unsuccess"})
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {authorization},
	}

	res, err := client.Do(req)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "request send error"})
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "readall error on res.Body"})
	}

	err = json.Unmarshal(body, &ln.ProfileInformation)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "unmarshal error"})
	}

	c.Status(200)
	return c.JSON(fiber.Map{"profile": ln.ProfileInformation})
}
