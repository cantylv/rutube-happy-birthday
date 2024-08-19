package entity

import "github.com/asaskevich/govalidator"

type SignUpForm struct {
	FullName string `json:"full_name" valid:"runelength(5|100)"`
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"password"`
	Birthday string `json:"birthday" valid:"date"`
}

func (d *SignUpForm) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

type SignInForm struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"password"`
}

func (d *SignInForm) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}
