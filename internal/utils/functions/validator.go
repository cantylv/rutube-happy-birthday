// Copyright © ivanlobanov. All rights reserved.
package functions

import (
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	valid "github.com/asaskevich/govalidator"
)

var (
	Email    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	Password = regexp.MustCompile(`^[A-Za-z\d!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]+$`)
	FullName = regexp.MustCompile(`^[A-Z][a-z]+ [A-Z][a-z]+ [A-Z][a-z]+$`)
)

// InitValidator
// Defines specific struct tags for validation.
func InitValidator() {
	valid.SetFieldsRequiredByDefault(true)
	// Custom validation tags
	govalidator.TagMap["user_email"] = func(email string) bool {
		emailLen := utf8.RuneCountInString(email)
		return emailLen >= 6 && emailLen <= 50 && Email.MatchString(email)
	}

	// user_phone_domain
	govalidator.TagMap["user_password"] = func(password string) bool {
		pwdLen := utf8.RuneCountInString(password)
		if pwdLen < 8 || pwdLen > 20 {
			return false
		}

		// Checking for the presence of at least one letter
		letterRegex := regexp.MustCompile(`[A-Za-z]`)
		if !letterRegex.MatchString(password) {
			return false
		}

		// Checking for the presence of at least one digit
		digitRegex := regexp.MustCompile(`\d`)
		if !digitRegex.MatchString(password) {
			return false
		}
		return Password.MatchString(password)
	}

	// user_email_domain
	govalidator.TagMap["user_fullname"] = func(fullname string) bool {
		nameLen := utf8.RuneCountInString(fullname)
		return nameLen >= 10 && nameLen <= 100 && FullName.MatchString(fullname)
	}

	// user_email_domain
	govalidator.TagMap["user_birthday"] = func(birthday string) bool {
		layout := "02.01.2006"
		_, err := time.Parse(layout, birthday)
		return err == nil
	}
}
