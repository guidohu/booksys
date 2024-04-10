package email

import (
	"fmt"
	"net/smtp"
	"server/database"
	"strings"
	"text/template"

	"golang.org/x/exp/slog"

	"embed"
)

//go:embed templates/*
var templatesFS embed.FS

type Client struct {
	Host     string
	Port     string
	Username string
	Password string
}

type Email struct {
	Recipient string
	Sender    string
	Message   string
}

func (c *Client) Send(e *Email) error {
	auth := smtp.PlainAuth("", c.Username, c.Password, c.Host)
	err := smtp.SendMail(fmt.Sprintf("%s:%s", c.Host, c.Port), auth, e.Sender, []string{e.Recipient}, []byte(e.Message))
	if err != nil {
		slog.Error("Email not sent", slog.String("to", e.Recipient), slog.String("error", err.Error()))
		return err
	}
	slog.Info("Email sent", slog.String("to", e.Recipient))
	return nil
}

func (c *Client) TokenResetMessage(user database.User, token string, sender string, subject string) (string, error) {
	var out strings.Builder
	tmplFile := "tokenReset.tmpl"
	tmpl := template.Must(template.New(tmplFile).ParseFS(templatesFS, "templates/tokenReset.tmpl"))
	vars := struct {
		FirstName string
		LastName  string
		Token     string
		Recipient string
		Sender    string
		Subject   string
	}{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Token:     token,
		Recipient: user.Email,
		Sender:    sender,
		Subject:   subject,
	}
	err := tmpl.Execute(&out, vars)
	if err != nil {
		slog.Error("Cannot generate output from template file", slog.String("error", err.Error()))
		return "", err
	}
	return out.String(), nil
}
