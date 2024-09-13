package email

import (
	"net/smtp"
)

// TODO: move to config email addr

// New send email message
func New(to string, subject string, body string) error {
	from := "testemail@example.com"
	password := "testEmailPasswd"

	smtpServer := "smtp.example.com:587"

	message := []byte("Subject: " + subject + "\r\n" + body)

	auth := smtp.PlainAuth("", from, password, "smtp.example.com")
	_ = smtp.SendMail(smtpServer, auth, from, []string{to}, message)
	// Handle error
	/*if err != nil {
		return err
	}*/

	return nil
}
