package golinkedin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Post struct {
	// Author URN for this content
	Author string `json:"author"`

	// The state of this content. PUBLISHED is the only accepted field during creation.
	LifeCycleState string `json:"lifecycleState"`

	// The content of post. For now you can just define text.
	SpecificContent SpecificContent `json:"specificContent"`

	// Visibility restrictions on content.
	Visibility Visibility `json:"visibility"`
}

type SpecificContent struct {
	ShareContent ShareContent `json:"com.linkedin.ugc.ShareContent"`
}

type ShareContent struct {
	ShareCommentary    ShareCommentary `json:"shareCommentary"`
	ShareMediaCategory string          `json:"shareMediaCategory"`
}

type ShareCommentary struct {
	Text string `json:"text"`
}

type Visibility struct {
	Code string `json:"com.linkedin.ugc.MemberNetworkVisibility"`
}

// JOB POSTING STRUCTS

type JobValue struct {
	JobPosting []JobPosting `json:"elements"`
}

type JobPosting struct {
	IntegrationContext      string   `json:"integrationContext"`
	CompanyApplyUrl         string   `json:"companyApplyUrl"`
	Description             string   `json:"description"`
	EmploymentStatus        string   `json:"employmentStatus"`
	ExternalJobPostingId    string   `json:"externalJobPostingId"`
	ListedAt                int      `json:"listedAt"`
	JobPostingOperationType string   `json:"jobPostingOperationType"`
	Title                   string   `json:"title"`
	Location                string   `json:"location"`
	WorkplaceTypes          []string `json:"workplaceTypes"`
}

/*
Have ProfileID, Share any post what you want.
Must post this route to json.

	{
		"text": <text>
	}
*/
func (ln *Linkedin) SharePost(c *fiber.Ctx) error {
	var dataBody map[string]string
	if err := c.BodyParser(&dataBody); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "body parser error"})
	}

	token := c.Cookies("linkedin_token")
	postUrl := fmt.Sprintf("https://api.linkedin.com/v2/ugcPosts?oauth2_access_token=%s", token)

	data := Post{
		Author:         "urn:li:person:" + ln.ProfileInformation.Id,
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

func (ln *Linkedin) ShareJOBPosting(c *fiber.Ctx) error {
	token := c.Cookies("linkedin_token")

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
