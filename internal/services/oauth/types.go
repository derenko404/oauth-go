package ouathservice

import "strconv"

type Profile interface {
	GetID() string
	GetEmail() string
	GetName() string
	GetAvatarURL() string
}

type ProfileImpl struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type GoogleProfile struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (g GoogleProfile) GetID() string {
	return g.ID
}

func (g GoogleProfile) GetEmail() string {
	return g.Email
}

func (g GoogleProfile) GetName() string {
	return g.Name
}

func (g GoogleProfile) GetAvatarURL() string {
	return g.Picture
}

type GithubProfile struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func (g GithubProfile) GetID() string {
	return strconv.Itoa(g.ID)
}

func (g GithubProfile) GetEmail() string {
	return g.Email
}

func (g GithubProfile) GetName() string {
	return g.Name
}

func (g GithubProfile) GetAvatarURL() string {
	return g.AvatarURL
}
