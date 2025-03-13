package user

import (
	"net/mail"
)

func isValidEmail(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}
