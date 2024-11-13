package services

import (
	"unicode"

	fails "github.com/BeatEcoprove/identityService/pkg/errors"
)

func containsNumber(s string) bool {
	for _, char := range s {
		if unicode.IsDigit(char) {
			return true
		}
	}

	return false
}

func containsCapitalLetter(s string) bool {
	for _, char := range s {
		if unicode.IsUpper(char) {
			return true
		}
	}

	return false
}

func containsNonCapitalLetter(s string) bool {
	for _, char := range s {
		if unicode.IsLower(char) {
			return true
		}
	}

	return false
}

func ValidatePassword(password string) error {
	if password == "" {
		return fails.PASSWORD_PROVIDE
	}

	if len(password) < 6 || len(password) > 16 {
		return fails.PASSWORD_BTW_6_16
	}

	if !containsNumber(password) {
		return fails.PASSWORD_MUST_CONTAIN_ONE_NUMBER
	}

	if !containsCapitalLetter(password) {
		return fails.PASSWORD_MUST_CONTAIN_ONE_CAPITAL
	}

	if !containsNonCapitalLetter(password) {
		return fails.PASSWORD_MUST_CONTAIN_NON_CAPITAL
	}

	return nil
}
