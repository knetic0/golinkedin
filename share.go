package golinkedin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

/*
Have ProfileID, Share any post what you want.
Must post this route to json.

	{
		"text": <text>
	}
*/
func (s *API) SharePost(c *fiber.Ctx) error {
	var dataBody map[string]string
	if err := c.BodyParser(&dataBody); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "body parser error"})
	}

	sess, err := AccessToken.Get(c)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "session get error"})
	}
	token := sess.Get("access_token")
	postUrl := fmt.Sprintf("https://api.linkedin.com/v2/ugcPosts?oauth2_access_token=%s", token)

	data := Post{
		Author:         "urn:li:person:" + s.ProfileInformation.Id,
		LifeCycleState: "PUBLISHED",
		SpecificContent: SpecificContent{
			ShareContent: ShareContent{
				ShareCommentary: ShareCommentary{
					Text: dataBody["text"],
				},
				ShareMediaCategory: "NONE",
			},
		},
		Visibility: Visibility{
			Code: "PUBLIC",
		},
	}

	post, err := postToJson(&data)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "post to json error"})
	}

	resp, err := http.Post(postUrl, "application/json", bytes.NewReader(post))
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "share post error"})
	}

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "read all error"})
	}

	c.Status(200)
	return c.JSON(fiber.Map{"message": bodyData})
}

func (s *API) ShareJOBPosting(c *fiber.Ctx) error {
	sess, err := AccessToken.Get(c)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "error on taking access token from cookie"})
	}
	token := sess.Get("access_token")

	authString := fmt.Sprintf("Bearer %s", token)
	jobPostingURL := "https://api.linkedin.com/v2/simpleJobPostings"
	client := http.Client{}

	data := JobValue{
		JobPosting: []JobPosting{
			{
				IntegrationContext:      "urn:li:organization:<organization-id>",
				CompanyApplyUrl:         "<company-url>",
				Description:             "We are looking for a passionate Software Engineer",
				EmploymentStatus:        "PART_TIME",
				ExternalJobPostingId:    "1234",
				ListedAt:                14400002023,
				JobPostingOperationType: "CREATE",
				Title:                   "Software Engineer",
				Location:                "Turkey",
				WorkplaceTypes:          []string{"remote"},
			},
		},
	}

	post, err := jobPostToJson(&data)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "post to json is not successfully"})
	}

	req, err := http.NewRequest(http.MethodPost, jobPostingURL, bytes.NewReader(post))
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "posting error"})
	}

	req.Header = http.Header{
		"Content-Type":              {"application/json"},
		"Authorization":             {authString},
		"X-Restli-Protocol-Version": {"2.0.0"},
		"X-Restli-Method":           {"batch_create"},
	}

	resp, err := client.Do(req)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "client do error"})
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "read all error"})
	}

	c.Status(200)
	return c.JSON(fiber.Map{"message": string(body)})
}

/*
Convert Post type to []byte type.
*/
func postToJson(post *Post) ([]byte, error) {
	data, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}

	return data, nil
}

/*
Convert Post type to []byte type.
*/
func jobPostToJson(post *JobValue) ([]byte, error) {
	data, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}

	return data, nil
}
