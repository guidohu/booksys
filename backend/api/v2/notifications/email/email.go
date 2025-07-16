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
	Subject   string
}

func NewClient(config database.EmailConfiguration) *Client {
	return &Client{
		Host:     config.Server,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
	}
}

func (c *Client) Send(e *Email) error {
	auth := smtp.PlainAuth("", c.Username, c.Password, c.Host)
	err := smtp.SendMail(fmt.Sprintf("%s:%s", c.Host, c.Port), auth, e.Sender, []string{e.Recipient}, []byte(e.Message))
	if err != nil {
		slog.Error("Email not sent", slog.String("to", e.Recipient), slog.String("error", err.Error()))
		return err
	}
	slog.Info("Email sent", slog.String("to", e.Recipient), slog.String("subject", e.Subject))
	return nil
}

func (c *Client) SendTokenResetMessage(user database.User, token string, senderAddress string, subject string) error {
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
		Sender:    senderAddress,
		Subject:   subject,
	}
	err := tmpl.Execute(&out, vars)
	if err != nil {
		slog.Error("Cannot generate output from template file", slog.String("error", err.Error()))
		return err
	}
	return c.Send(&Email{
		Recipient: user.Email,
		Sender:    senderAddress,
		Message:   out.String(),
		Subject:   subject,
	})
}

func (c *Client) SendUserAddedToSessionMessage(user database.User, session database.Session, senderAddress string, subject string, URL string) error {
	var out strings.Builder
	tmplFile := "userAddedToSession.tmpl"
	tmpl := template.Must(template.New(tmplFile).ParseFS(templatesFS, "templates/userAddedToSession.tmpl"))
	vars := struct {
		FirstName string
		LastName  string
		Date      string
		Start     string
		End       string
		Recipient string
		Sender    string
		Subject   string
		URL       string
	}{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Date:      session.StartTime.Format("Monday 02.01.2006"),
		Start:     session.StartTime.Format("15:04"),
		End:       session.EndTime.Format("15:04"),
		Recipient: user.Email,
		Sender:    senderAddress,
		Subject:   subject,
		URL:       URL,
	}
	err := tmpl.Execute(&out, vars)
	if err != nil {
		slog.Error("Cannot generate output from template file", slog.String("error", err.Error()))
		return err
	}
	return c.Send(&Email{
		Recipient: user.Email,
		Sender:    senderAddress,
		Message:   out.String(),
		Subject:   subject,
	})
}

func (c *Client) SendUserRemovedFromSessionMessage(user database.User, session database.Session, senderAddress string, subject string, URL string) error {
	var out strings.Builder
	tmplFile := "userRemovedFromSession.tmpl"
	tmpl := template.Must(template.New(tmplFile).ParseFS(templatesFS, "templates/userRemovedFromSession.tmpl"))
	vars := struct {
		FirstName string
		LastName  string
		Date      string
		Start     string
		End       string
		Recipient string
		Sender    string
		Subject   string
		URL       string
	}{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Date:      session.StartTime.Format("Monday 02.01.2006"),
		Start:     session.StartTime.Format("15:04"),
		End:       session.EndTime.Format("15:04"),
		Recipient: user.Email,
		Sender:    senderAddress,
		Subject:   subject,
		URL:       URL,
	}
	err := tmpl.Execute(&out, vars)
	if err != nil {
		slog.Error("Cannot generate output from template file", slog.String("error", err.Error()))
		return err
	}
	return c.Send(&Email{
		Recipient: user.Email,
		Sender:    senderAddress,
		Message:   out.String(),
		Subject:   subject,
	})
}
