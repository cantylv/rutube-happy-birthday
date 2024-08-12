package entity

type Subscription struct {
	EmployeeId uint32 `json:"employee_id" valid:"-" bson:"employee_id"`
	Interval   uint16 `json:"interval" valid:"-" bson:"interval"`
	IsFollowed bool   `json:"is_followed" valid:"-" bson:"is_followed"`
}
