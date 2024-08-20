package entity

import "github.com/asaskevich/govalidator"

type Subscription struct {
	EmployeeId       string `json:"employee_id" valid:"-" bson:"employee_id"`
	EmployeeFullName string `json:"employee_full_name" valid:"user_fullname" bson:"employee_full_name"`
	EmployeeBirthday string `json:"employee_birthday" valid:"user_birthday" bson:"employee_birthday"`
	EmployeeEmail    string `json:"employee_email" valid:"user_email" bson:"employee_email"`
	Interval         uint16 `json:"interval" valid:"-" bson:"interval"`
	IsFollowed       bool   `json:"is_followed" valid:"-" bson:"is_followed"`
}

func (d *Subscription) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

type SubProps struct {
	IdFollower string
	IdEmployee string
}

type SetUpIntervalProps struct {
	Ids         SubProps
	NewInterval uint16
}
