package services

import (
	"fmt"

	"github.com/FebryanHernanda/BE-EventOrganizer/internal/mailer/provider"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/mailer/templates"
	"github.com/sirupsen/logrus"
)

type MailerService struct {
	provider *provider.SMTPProvider
}

func NewMailerService(provider *provider.SMTPProvider) *MailerService {
	return &MailerService{provider}
}

func (s *MailerService) SendActivationEmail(name, email, registeredDate, token string) error {
	link := fmt.Sprintf("http://localhost:8080/auth/activate?token=%s", token)

	body := templates.ActivationTemplate(name, email, registeredDate, link)

	if err := s.provider.Send(email, "Activate Your Account", body); err != nil {
		logrus.WithError(err).Error("Failed sending activation email")
		return err
	}

	logrus.Info("Activation email sent successfully")
	return nil
}
