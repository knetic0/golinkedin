package golinkedin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

/*
After than Callback, you take Profile information with this function.
Create route for this.
*/
func (s *API) Profile(c *fiber.Ctx) error {
	profileUrl := "https://api.linkedin.com/v2/me"

	sess, err := AccessToken.Get(c)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "session get error"})
	}

	token := sess.Get("access_token")
	authHeader := fmt.Sprintf("Bearer %s", token)

	client := http.Client{}

	req, err := http.NewRequest("GET", profileUrl, nil)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "new request created unsuccess"})
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {authHeader},
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

	err = json.Unmarshal(body, &s.ProfileInformation)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "unmarshal error"})
	}

	c.Status(200)
	return c.JSON(fiber.Map{"profile": s.ProfileInformation})
}
