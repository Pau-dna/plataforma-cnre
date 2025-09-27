package validators

import (
	"errors"
	"net/mail"
	"strings"
)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if len(email) == 0 {
		return errors.New("el email no puede estar vacío")
	}

	if len(email) > 254 {
		return errors.New("el email no puede exceder 254 caracteres")
	}

	if !strings.Contains(email, "@") {
		return errors.New("formato de email inválido")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("formato de email inválido")
	}
	return nil
}
