package services

import (
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

func (s *EmailServiceResentImpl) SendAccountConfirmation(name string, email string, otp string) error {

	params := &resend.SendEmailRequest{
		From:    common.NOREPLY_EMAIL,
		To:      []string{email},
		Subject: "Confirm Your Account!",
		Html:    "<p>Congrats on sending your <strong>" + otp + "</strong>!</p>",
	}

	_, err := s.resendClient.Emails.Send(params)

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
