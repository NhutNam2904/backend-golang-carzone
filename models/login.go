package models

import "errors"

type Credentials struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func ValidateCredentials(credential Credentials) error {

	if err := ValidateUsername(credential.UserName); err != nil {
		return err
	}
	if err := ValidatePassword(credential.Password); err != nil {
		return err
	}

	return nil

}

func ValidateUsername(username string) error {

	if username == "" {
		return errors.New("Please enter username")
	}
	return nil

}
func ValidatePassword(password string) error {

	if password == "" {
		return errors.New("Please enter password")
	}
	return nil

}
