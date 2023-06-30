package golinkedin

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

	// ProfileInformation inherit on API Struct.
	ProfileInformation ProfileInformation `json:"profile_information"`
}

type ProfileInformation struct {
	// ProfileID represents the ID every linkedin profile has.
	Id string `json:"id"`

	// User's FirstName
	FirstName string `json:"localizedFirstName"`

	// User's LastName
	LastName string `json:"localizedLastName"`
}

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
