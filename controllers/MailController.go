package controllers

import (
	"fmt"
	"net/mail"
	"net/smtp"

	"github.com/scorredoira/email"
)

type Request struct {
	from    string
	bcc     *string
	to      []string
	subject string
	body    string
}

type FileRequest struct {
	to      []string
	subject string
	body    string
	file    string
}

func NewRequest(to []string, subject string, body string, bcc *string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
		bcc:     bcc,
	}
}

func NewFileRequest(to []string, subject string, body string, file string) *FileRequest {
	return &FileRequest{
		to:      to,
		subject: subject,
		body:    body,
		file:    file,
	}
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	from := "From: Wastecontrol <support@wastecontrol.world>\n"
	subject := "Subject: " + r.subject + "\n"
	msg := []byte(from + subject + mime + "\n" + r.body)
	addr := "smtp.gigahost.dk:587"
	auth := smtp.PlainAuth("", "support@wastecontrol.world", "GolfLars123", "smtp.gigahost.dk")

	if err := smtp.SendMail(addr, auth, "support@wastecontrol.world", r.to, msg); err != nil {
		fmt.Println(err)
		return false, err
	}

	return true, nil
}

func (fr *FileRequest) SendEmailWithFile() {
	auth := smtp.PlainAuth("", "support@wastecontrol.world", "GolfLars123", "smtp.gigahost.dk")
	from := "support@wastecontrol.world"
	to := fr.to
	subject := "Subject: " + fr.subject + "\n"
	msg := fr.body

	emailContent := email.NewMessage(subject, msg)

	emailContent.From = mail.Address{Name: "Wastecontrol", Address: from}
	emailContent.To = to

	err := emailContent.Attach(fr.file)
	if err != nil {
		fmt.Println(err)
	}

	err = email.Send("smtp.gigahost.dk:587", auth, emailContent)
	if err != nil {
		fmt.Println(err)
	}
}
