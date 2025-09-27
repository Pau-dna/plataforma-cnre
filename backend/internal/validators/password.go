package validators

import (
	"errors"
	"strings"
)

func ValidatePassword(password string) error {

	cleanedPassword := strings.TrimSpace(password)

	if cleanedPassword == "" {
		return errors.New("la contraseña no puede estar vacía")
	}

	if len(cleanedPassword) < 8 {
		return errors.New("la contraseña debe tener al menos 8 caracteres")
	}

	if len(cleanedPassword) > 100 {
		return errors.New("la contraseña no puede exceder 100 caracteres")
	}

	// Check insecure passwords
	insecurePasswords := []string{"12345678", "password"}
	for _, insecurePassword := range insecurePasswords {
		if cleanedPassword == insecurePassword {
			return errors.New("la contraseña es muy débil")
		}
	}

	return nil
}
