package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

type Email struct {
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

func NewEmail(name, user, password, host string) *Email {
	hp := strings.Split(host, ":")
	if len(hp) < 2 {
		host = strings.Join([]string{hp[0], "587"}, ":")
	}
	return &Email{
		Name:     name,
		User:     user,
		Password: password,
		Host:     host,
	}
}

func (e *Email) Send(ctx context.Context, tos []string, mailtype, subject, body string) error {
	auth := smtp.PlainAuth("", e.User, e.Password, e.Host)
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	to := strings.Join(tos, ";")
	from := e.Name
	if from == "" {
		from = e.User
	}
	msg := fmt.Sprintf("To:%s\r\nFrom:%s<%s>\r\nSubject:%s\r\n%s\r\n\r\n%s",
		to, from, e.User, subject, contentType, body)
	//return smtp.SendMail(e.Host, auth, e.User, tos, []byte(msg))
	conn, err := tls.Dial("tcp", e.Host, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return err
	}
	co, err := smtp.NewClient(conn, e.Host)
	if err != nil {
		return err
	}
	defer func() {
		_ = co.Close()
	}()
	if ok, _ := co.Extension("AUTH"); ok {
		if err = co.Auth(auth); err != nil {
			return err
		}
	}
	if err = co.Mail(e.User); err != nil {
		return err
	}
	for _, addr := range tos {
		if err = co.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := co.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return co.Quit()
}
