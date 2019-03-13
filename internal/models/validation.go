package models

import (
	"escapade/internal/misc"
)

func ValidateStrings(text ...string) (nullstring bool) {
	nullstring = false
	for _, str := range text {
		if len(str) == 0 {
			nullstring = true
			break
		}
	}
	return
}

func ValidateString(str string) bool {
	return len(str) > 0
}

func ValidatePlayerName(str string) bool {
	return ValidateString(str)
}

func ValidateEmail(str string) bool {
	return ValidateString(str)
}

func ValidatePassword(str string) bool {
	return ValidateString(str)
}

func ValidateCookieName(str string) bool {
	return str == misc.NameCookie
}

func ValidateCookieValue(str string) bool {
	return misc.LengthCookie == len(str)
}
