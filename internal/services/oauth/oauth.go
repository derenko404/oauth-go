package ouathservice

import (
	"context"
	"encoding/json"
	"fmt"
	"go-auth/internal/types"
	"io"
	"net/http"
	"slices"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var SupportedProviders = []string{"google", "github"}
var ProfileURLs = map[string]string{
	"github": "https://api.github.com/user",
	"google": "https://www.googleapis.com/oauth2/v2/userinfo",
}

type OAuthService interface {
	IsSupported(provider string) error
	GetSignInUrl(provider string) (string, error)
	GetProfile(ctx context.Context, provider string, code string)
}

type OAuth struct {
	config *types.AppConfig
}

func New(config *types.AppConfig) *OAuth {
	return &OAuth{
		config: config,
	}
}

func getProfile[T Profile](provider string, client *http.Client) (*ProfileImpl, error) {
	url, ok := ProfileURLs[provider]

	if !ok {
		return nil, fmt.Errorf("no url for selected provider %s", url)
	}

	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var info T
	err = json.Unmarshal(body, &info)

	if err != nil {
		return nil, err
	}

	return &ProfileImpl{
		ID:        info.GetID(),
		Email:     info.GetEmail(),
		Name:      info.GetName(),
		AvatarURL: info.GetAvatarURL(),
	}, nil
}

func (oauth *OAuth) getConfig(provider string) (*oauth2.Config, error) {
	err := oauth.IsSupported(provider)

	if err != nil {
		return nil, err
	}

	if provider == "google" {
		return &oauth2.Config{
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
			ClientID:     oauth.config.GoogleClientId,
			ClientSecret: oauth.config.GoogleClientSecret,
			RedirectURL:  oauth.config.GoogleRedirectURL,
		}, nil
	}

	if provider == "github" {
		return &oauth2.Config{
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     github.Endpoint,
			ClientID:     oauth.config.GithubClientId,
			ClientSecret: oauth.config.GithubClientSecret,
			RedirectURL:  oauth.config.GithubRedirectURL,
		}, nil
	}

	return nil, fmt.Errorf("unsuported provider")
}

func (oauth *OAuth) IsSupported(provider string) error {
	if provider == "" {
		return fmt.Errorf("provider is required")
	}

	if !slices.Contains(SupportedProviders, provider) {
		return fmt.Errorf("unsupported oauth provider")
	}

	return nil
}

func (oauth *OAuth) GetSignInUrl(provider string, state string) (string, error) {
	err := oauth.IsSupported(provider)

	if err != nil {
		return "", err
	}

	oAuthConfig, err := oauth.getConfig(provider)

	if err != nil {
		return "", fmt.Errorf("cannot get oauth config %w", err)
	}

	url := oAuthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)

	return url, nil
}

func (oauth *OAuth) GetProfile(ctx context.Context, provider string, code string) (*ProfileImpl, error) {
	err := oauth.IsSupported(provider)

	if err != nil {
		return nil, err
	}

	oAuthConfig, err := oauth.getConfig(provider)

	if err != nil {
		return nil, fmt.Errorf("cannot get oauth config %w", err)
	}

	tokens, err := oAuthConfig.Exchange(ctx, code)

	if err != nil {
		return nil, fmt.Errorf("cannot exchange tokens %w", err)
	}

	client := oAuthConfig.Client(ctx, tokens)

	if provider == "google" {
		user, err := getProfile[GoogleProfile](provider, client)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	if provider == "github" {
		user, err := getProfile[GithubProfile](provider, client)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, fmt.Errorf("unsuported profile")
}
