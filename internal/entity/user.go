package entity

import "github.com/asaskevich/govalidator"

type User struct {
	Id       string         `json:"id" valid:"-"`
	FullName string         `json:"full_name" valid:"user_fullname"`
	Birthday string         `json:"birthday" valid:"user_birthday"`
	Email    string         `json:"email" valid:"user_email"`
	Password string         `json:"password,omitempty" valid:"password"`
	Subs     []Subscription `json:"subs,omitempty" valid:"-"`
}

type UserWithoutPassword struct {
	Id       string         `json:"id"`
	FullName string         `json:"full_name"`
	Birthday string         `json:"birthday"`
	Email    string         `json:"email"`
	Subs     []Subscription `json:"subs"`
}

func (d *User) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

type UserUpdate struct {
	FullName string `json:"full_name" valid:"user_fullname"`
	Birthday string `json:"birthday" valid:"user_birthday"`
	Email    string `json:"email" valid:"user_email"`
	Password string `json:"password,omitempty" valid:"user_password"`
}

func (d *UserUpdate) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}
