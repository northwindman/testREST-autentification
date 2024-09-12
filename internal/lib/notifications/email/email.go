package email

import (
	"errors"
	"net/smtp"
)

var (
	ResOk = errors.New("email send complete")
)

// New send email message
func New(to string, subject string, body string) error {
	from := "testemail@example.com"
	password := "testEmailPasswd"

	smtpServer := "smtp.example.com:587"

	message := []byte("Subject: " + subject + "\r\n" + body)

	auth := smtp.PlainAuth("", from, password, "smtp.example.com")
	smtp.SendMail(smtpServer, auth, from, []string{to}, message)
	// Handle error
	/*if err != nil {
		return err
	}*/

	return nil
}
