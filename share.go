package golinkedin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

var (
	PostURL = "https://api.linkedin.com/v2/ugcPosts?oauth2_access_token="
)

/*
Have ProfileID, Share any post what you want.
*/
func (ln *Linkedin) SharePost(token, id, text string) error {
	postUrl := fmt.Sprintf("%s%s", PostURL, token)

	data := Post{
		Author:         "urn:li:person:" + id,
		LifeCycleState: "PUBLISHED",
		SpecificContent: SpecificContent{
			ShareContent: ShareContent{
				ShareCommentary: ShareCommentary{
					Text: text,
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
		return fmt.Errorf("error occurred = %s", err.Error())
	}

	resp, err := http.Post(postUrl, "application/json", bytes.NewReader(post))
	if err != nil {
		return fmt.Errorf("error occurred = %s", err.Error())
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error occurred = %s", err.Error())
	}

	return nil
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
