package emails

import (
	"encoding/json"
	"fmt"
	"net/smtp"
)

//EmailAddress contains all data related to an Email address for Server Side Usage
type EmailAddress struct {
	Address  string
	Port     int
	Password string
	Host     string
}



//SendEmailHTML Sends an email to single user
func (emailAddress *EmailAddress) SendEmailHTML(to string, emailSubject string, htmlCode string) (bool, error) {

	var emailAuth smtp.Auth
	emailAuth = smtp.PlainAuth("", emailAddress.Address, emailAddress.Password, emailAddress.Host)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + emailSubject + "!\n"
	msg := []byte(subject + mime + htmlCode)
	addr := fmt.Sprintf("%s:%d", emailAddress.Host, emailAddress.Port)

	if err := smtp.SendMail(addr, emailAuth, emailAddress.Address, []string{to}, msg); err != nil {
		return false, err
	}
	return true, nil
}

//SendEmail Sends an email to single user
func (emailAddress *EmailAddress) SendEmail(to string, emailSubject string, emailBody string) (bool, error) {

	bytes, err := json.Marshal(*emailAddress)
	if err != nil {
		return false, err
	}

	fmt.Println(string(bytes))

	var emailAuth smtp.Auth
	emailAuth = smtp.PlainAuth("", emailAddress.Address, emailAddress.Password, emailAddress.Host)

	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + emailSubject + "!\n"
	msg := []byte(subject + mime + "\n" + emailBody)
	addr := fmt.Sprintf("%s:%d", emailAddress.Host, emailAddress.Port)

	fmt.Println(addr)

	if err := smtp.SendMail(addr, emailAuth, emailAddress.Address, []string{to}, msg); err != nil {
		return false, err
	}
	return true, nil
}

