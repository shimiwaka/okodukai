package main

import (
	"fmt"
	"os"
    "net/smtp"
	"net/http"
    "strings"

	"github.com/shimiwaka/okodukai/schema"
	"github.com/shimiwaka/okodukai/connector"
	"gopkg.in/yaml.v2"
)

type MailSettings struct {
	Hostname string `yaml:"mail_host"`
	MailAddress string `yaml:"mail_address"`
	MailPassword string `yaml:"mail_password"`
}

func forgetHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()

	e := r.ParseForm()
	if e != nil {
		panic("error: parse error occured.")
	}

	email := r.Form.Get("email")

	boards := []schema.Board{}

	db.Where("owner = ?", email).Find(&boards)

	mailSettings := MailSettings{}
	b, _ := os.ReadFile("mail.yaml")
	yaml.Unmarshal(b, &mailSettings)

    from := mailSettings.MailAddress
    recipients := []string{email}
    subject := "okodukai information"
    body := "The board assigned to your email address is below.\n\n"

	for i := 0; i < len(boards); i++ {
		body += boards[i].Token + "\n"
	}

    auth := smtp.CRAMMD5Auth(mailSettings.MailAddress, mailSettings.MailPassword)
    msg := []byte(strings.ReplaceAll(fmt.Sprintf("To: %s\nSubject: %s\n\n%s", strings.Join(recipients, ","), subject, body), "\n", "\r\n"))

	err := smtp.SendMail(fmt.Sprintf("%s:%d", mailSettings.Hostname, 587), auth, from, recipients, msg);
    if err != nil {
        panic("error : failed to send mail.")
    }
	fmt.Fprintln(w, "{\"success\": true}")

	db.Close()
}
