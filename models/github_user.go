package models

type User struct {
	Login           string `json:"login"`
	ID              int    `json:"id"`
	NodeID          string `json:"node_id"`
	AvatarURL       string `json:"avatar_url"`
	GravatarID      string `json:"gravatar_id"`
	HTMLURL         string `json:"html_url"`
	Type            string `json:"type"`
	SiteAdmin       bool   `json:"site_admin"`
	Name            string `json:"name"`
	Company         string `json:"company"`
	Blog            string `json:"blog"`
	Location        string `json:"location"`
	Email           string `json:"email"`
	EmailRepo       string `json:"emailRepo"`
	Bio             string `json:"bio"`
	TwitterUsername string `json:"twitter_username"`
	PublicRepos     int    `json:"public_repos"`
	PublicGists     int    `json:"public_gists"`
	Followers       int    `json:"followers"`
	Following       int    `json:"following"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type UserRepo struct {
	Fork       bool   `json:"fork"`
	CommitsURL string `json:"commits_url"`
	CreatedAt  string `json:"created_at"`
}
