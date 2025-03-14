package email

import "github.com/sendgrid/sendgrid-go/helpers/mail"

type SendgridPayload struct {
	From      *mail.Email
	To        *mail.Email
	PlainText string
	HTMLText  string
	Subject   string
}
