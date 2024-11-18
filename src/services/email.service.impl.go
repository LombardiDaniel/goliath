package services

import (
	"bytes"
	"log/slog"
	"net/url"
	"path/filepath"
	"text/template"

	"github.com/LombardiDaniel/go-gin-template/common"
	"github.com/resendlabs/resend-go"
)

type EmailServiceResentImpl struct {
	resendClient                *resend.Client
	accountConfirmationTemplate *template.Template
	accountCreationTemplate     *template.Template
}

func NewEmailServiceResentImpl(resendApiKey string, templatesDir string) EmailService {
	return &EmailServiceResentImpl{
		resendClient:                resend.NewClient(resendApiKey),
		accountConfirmationTemplate: common.LoadHTMLTemplate(filepath.Join(templatesDir, "account-confirmation.html")),
		accountCreationTemplate:     common.LoadHTMLTemplate(filepath.Join(templatesDir, "account-created.html")),
	}
}

type htmlConfirmationVars struct {
	Name            string
	ConfirmationUrl string
}

func (s *EmailServiceResentImpl) SendAccountConfirmation(name string, email string, otp string) error {

	body := new(bytes.Buffer)
	confirmUrl, err := url.JoinPath(common.API_HOST_URL, "/v1/users/confirm")
	if err != nil {
		return err
	}
	err = s.accountConfirmationTemplate.Execute(body, htmlConfirmationVars{
		Name:            name,
		ConfirmationUrl: confirmUrl + "?otp=" + otp,
	})
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	params := &resend.SendEmailRequest{
		From:    common.NOREPLY_EMAIL,
		To:      []string{email},
		Subject: "Confirm Your Account!",
		Html:    body.String(),
	}

	_, err = s.resendClient.Emails.Send(params)

	return err
}

func (s *EmailServiceResentImpl) SendAccountCreated(name string, email string) error {
	params := &resend.SendEmailRequest{
		From:    common.NOREPLY_EMAIL,
		To:      []string{email},
		Subject: "Account Created!",
		Html:    "<p>Congrats on sending your <strong>" + name + "</strong>!</p>",
	}

	_, err := s.resendClient.Emails.Send(params)

	return err
}
