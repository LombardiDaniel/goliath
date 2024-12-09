package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/LombardiDaniel/gopherbase/common"
	"golang.org/x/oauth2"
)

const (
	githubUserInfoRetrievalUrl string = "https://api.github.com/user"
)

type GithubProvider struct {
	Config *oauth2.Config

	// private
	authUrl string
}

func NewGithubProvider(conf *oauth2.Config) Provider {
	return &GithubProvider{
		Config:  conf,
		authUrl: conf.AuthCodeURL(""),
	}
}

func (p *GithubProvider) GetAuthUrl() string {
	return p.authUrl
}

func (p *GithubProvider) Auth(ctx context.Context, code string) (*User, error) {
	token, err := p.Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, githubUserInfoRetrievalUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("body: %v\n", string(body))

	usrSchema := githubOauthSchema{}
	err = json.Unmarshal(body, &usrSchema)
	if err != nil {
		return nil, err
	}

	first, last := common.SplitName(usrSchema.Name)

	user := User{
		Email:        usrSchema.Email,
		FirstName:    first,
		LastName:     last,
		PictureUrl:   &usrSchema.AvatarURL,
		Provider:     GITHUB_PROVIDER,
		RefreshToken: token.RefreshToken,
	}

	return &user, nil
}
