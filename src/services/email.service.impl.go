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

type EmailServiceResendImpl struct {
	resendClient                *resend.Client
	accountConfirmationTemplate *template.Template
	accountCreationTemplate     *template.Template
	organizationInviteTemplate  *template.Template
}

func NewEmailServiceResendImpl(resendApiKey string, templatesDir string) EmailService {
	return &EmailServiceResendImpl{
		resendClient:                resend.NewClient(resendApiKey),
		accountConfirmationTemplate: common.LoadHTMLTemplate(filepath.Join(templatesDir, "account-confirmation.html")),
		accountCreationTemplate:     common.LoadHTMLTemplate(filepath.Join(templatesDir, "account-created.html")),
		organizationInviteTemplate:  common.LoadHTMLTemplate(filepath.Join(templatesDir, "organization-invite.html")),
	}
}

type htmlConfirmationVars struct {
	Name            string
	ConfirmationUrl string
}

func (s *EmailServiceResendImpl) SendAccountConfirmation(name string, email string, otp string) error {

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

func (s *EmailServiceResendImpl) SendAccountCreated(name string, email string) error {
	params := &resend.SendEmailRequest{
		From:    common.NOREPLY_EMAIL,
		To:      []string{email},
		Subject: "Account Created!",
		Html:    "<p>Congrats on sending your <strong>" + name + "</strong>!</p>",
	}

	_, err := s.resendClient.Emails.Send(params)

	return err
}

func (s *EmailServiceResendImpl) SendOrganizationInvite(name string, email string, orgId string, orgName string) error {
	params := &resend.SendEmailRequest{
		From:    common.NOREPLY_EMAIL,
		To:      []string{email},
		Subject: "Organization Invite",
		Html:    "<p>Congrats on sending your <strong>" + name + orgId + orgName + "</strong>!</p>",
	}

	_, err := s.resendClient.Emails.Send(params)

	return err
}
