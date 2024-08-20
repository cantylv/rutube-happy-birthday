package entity

import "github.com/asaskevich/govalidator"

type SignUpForm struct {
	FullName string `json:"full_name" valid:"user_fullname"`
	Birthday string `json:"birthday" valid:"user_birthday"`
	Email    string `json:"email" valid:"user_email"`
	Password string `json:"password,omitempty" valid:"user_password"`
}

func (d *SignUpForm) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

type SignInForm struct {
	Email    string `json:"email" valid:"user_email"`
	Password string `json:"password,omitempty" valid:"user_password"`
}

func (d *SignInForm) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}
