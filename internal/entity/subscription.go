package entity

type Subscription struct {
	EmployeeId string `json:"employee_id" valid:"-" bson:"employee_id"`
	Interval   uint16 `json:"interval" valid:"-" bson:"interval"`
	IsFollowed bool   `json:"is_followed" valid:"-" bson:"is_followed"`
}

type SubProps struct {
	IdFollower string
	IdEmployee string
}

type SetUpIntervalProps struct {
	Ids         SubProps
	NewInterval uint16
}
