package entity

import "github.com/asaskevich/govalidator"

type User struct {
	Id       string         `json:"id" valid:"-"`
	FullName string         `json:"full_name" valid:"runelength(5|100)"`
	Birthday string         `json:"birthday" valid:"date"`
	Email    string         `json:"email" valid:"email"`
	Password string         `json:"password,omitempty" valid:"password"`
	Subs     []Subscription `json:"subs,omitempty" valid:"-"`
}

func (d *User) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

type UserUpdate struct {
	FullName string `json:"full_name" valid:"runelength(5|100)"`
	Birthday string `json:"birthday" valid:"date"`
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"password"`
}

func (d *UserUpdate) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}
