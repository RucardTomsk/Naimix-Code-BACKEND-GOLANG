package mail

import (
	"bytes"
	"fmt"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/common"
	"go.uber.org/zap"
	"net/mail"
	"net/smtp"
)

type MailService struct {
	smtpConfig common.SmtpConfig
	logger     *zap.Logger
}

func NewMailService(smtpConfig common.SmtpConfig, logger *zap.Logger) *MailService {
	return &MailService{
		smtpConfig: smtpConfig,
		logger:     logger,
	}
}

func (m *MailService) auth() smtp.Auth {
	return smtp.PlainAuth("", m.smtpConfig.User, m.smtpConfig.Password, m.smtpConfig.Host)
}

func (m *MailService) SendMessage(recipientAdder string, subject string, message string) error {
	from := mail.Address{Name: "Promitent", Address: m.smtpConfig.User}
	recipient := mail.Address{Name: "Promitent-юрист", Address: recipientAdder}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = recipient.String()
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	var msg bytes.Buffer
	for k, v := range header {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	msg.WriteString("\r\n" + message)

	if err := smtp.SendMail(m.smtpConfig.Host+":"+m.smtpConfig.Port, m.auth(), from.Address, []string{recipient.Address}, msg.Bytes()); err != nil {
		return err
	}

	return nil
}
