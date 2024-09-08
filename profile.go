package golinkedin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ProfileInformation struct {
	// ProfileID represents the ID every linkedin profile has.
	Id string `json:"id"`

	// User's FirstName
	FirstName string `json:"localizedFirstName"`

	// User's LastName
	LastName string `json:"localizedLastName"`
}

var (
	ProfileURL = "https://api.linkedin.com/v2/me"
)

/*
After than Callback, you take Profile information with this function.
Create route for this.
*/
func (ln *Linkedin) Profile(token string) (*ProfileInformation, error) {
	authorization := fmt.Sprintf("Bearer %s", token)

	client := http.Client{}

	req, _ := http.NewRequest("GET", ProfileURL, nil)

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {authorization},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request error: %s", err.Error())
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read body error: %s", err.Error())
	}

	var profile ProfileInformation

	_ = json.Unmarshal(body, &profile)

	return &profile, nil
}
