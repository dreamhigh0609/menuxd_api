package email

import (
	"bytes"
	"errors"
	"html/template"
	"os"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

// Errors
var (
	ErrImpossibleReadFile  = errors.New("impossible to read the template file")
	ErrImpossibleParseFile = errors.New("impossible to parse the template file")
)

var (
	// address Email address of the admin
	address  = os.Getenv("XD_EMAIL")
	password = os.Getenv("XD_PASSWORD")

	basePath = os.Getenv("XD_BASE_PATH")
)

// Email contains data to build and send emails
type Email struct {
	From         string `json:"from"`
	UserID       uint   `json:"user_id,omitempty"`
	Addr         string `json:"addr"`
	Subject      string `json:"subject"`
	Body         string `json:"body"`
	TempPassword string `json:"temp_password"`
	BaseURL      string `json:"base_url"`
	Template     string `json:"template"`
}

// smtpServer data to smtp server
type smtpServer struct {
	host string
	port string
}

// serverName URI to smtp server
func (s *smtpServer) serverName() string {
	return s.host + ":" + s.port
}

// buildMessage build a email body
func (e *Email) buildMessage() (string, error) {
	t, err := template.ParseFiles(e.Template)
	if err != nil {
		return "", ErrImpossibleReadFile
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, e)
	if err != nil {
		return "", ErrImpossibleParseFile
	}

	return buf.String(), nil
}

// NewUser returns a Email with settings to new user
func NewUser(uid uint, addr string, tpass string) Email {
	url := os.Getenv("BASE_URL_CLIENT")
	e := Email{
		From:         address,
		Subject:      "Bienvenido a MenuXD",
		BaseURL:      url,
		UserID:       uid,
		Addr:         addr,
		Body:         "¡Gracias por contratar nuestros servicios! Estas a un paso de volver tu menú tradicional a Formato Digital. ¡Si tienes dudas o consultas no dudes en escribirnos!",
		TempPassword: tpass,
		Template:     basePath + "template/email/email.html",
	}
	return e
}

// ChangePassword returns a Email with settings to change password
func ChangePassword(addr string, tpass string) Email {
	url := os.Getenv("BASE_URL_CLIENT")
	e := Email{
		From:         address,
		Subject:      "Contraseña temporal",
		BaseURL:      url,
		Addr:         addr,
		Body:         "¡Si tienes dudas o consultas no dudes en escribirnos!",
		TempPassword: tpass,
		Template:     basePath + "template/email/email.html",
	}
	return e
}

// Send send a Email
func (e Email) Send() error {
	message, err := e.buildMessage()
	if err != nil {
		return err
	}
	m := mailjet.NewMailjetClient(
		os.Getenv("MJ_APIKEY_PUBLIC"),
		os.Getenv("MJ_APIKEY_PRIVATE"),
	)

	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: address,
				Name:  "MenuXD",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: e.Addr,
				},
			},
			Subject:  e.Subject,
			HTMLPart: message,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}

	_, err = m.SendMailV31(&messages)
	if err != nil {
		return err
	}
	return nil
}
