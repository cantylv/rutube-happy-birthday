// Copyright © ivanlobanov. All rights reserved.
package functions

import (
	"regexp"

	valid "github.com/asaskevich/govalidator"
)

func InitValidator() {
	valid.SetFieldsRequiredByDefault(true)
	// Custom validation tags
	valid.ParamTagRegexMap["full_name"] = regexp.MustCompile(`^[A-Z][a-z]+ [A-Z][a-z]+ [A-Z][a-z]+$`)
	valid.ParamTagRegexMap["date"] = regexp.MustCompile(`^(19|20)\d{2}\.(0[1-9]|1[0-2])\.(0[1-9]|[1-2]\d|3[0-1])$`)
}
