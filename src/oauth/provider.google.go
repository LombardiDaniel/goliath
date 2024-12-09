package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

const (
	googleUserInfoRetrievalUrl string = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

type GoogleProvider struct {
	Config *oauth2.Config

	// private
	authUrl string
}

func NewGoogleProvider(conf *oauth2.Config) Provider {
	return &GoogleProvider{
		Config:  conf,
		authUrl: conf.AuthCodeURL(""),
	}
}

func (p *GoogleProvider) GetAuthUrl() string {
	return p.authUrl
}

func (p *GoogleProvider) Auth(ctx context.Context, code string) (*User, error) {
	token, err := p.Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(googleUserInfoRetrievalUrl + token.AccessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	usrSchema := googleOauthSchema{}
	err = json.Unmarshal(body, &usrSchema)
	if err != nil {
		return nil, err
	}

	if !usrSchema.VerifiedEmail {
		return nil, errors.New("user's email is not verified")
	}

	user := User{
		Email:        usrSchema.Email,
		FirstName:    usrSchema.GivenName,
		LastName:     usrSchema.FamilyName,
		PictureUrl:   &usrSchema.Picture,
		Provider:     GOOGLE_PROVIDER,
		RefreshToken: token.RefreshToken,
	}

	return &user, nil
}
