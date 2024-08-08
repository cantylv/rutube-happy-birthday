package entity

type Subscription struct {
	FollowerId uint32 `json:"follower_id" valid:"-"`
	EmployeeId uint32 `json:"employee_id" valid:"-"`
	Interval   uint16 `json:"interval" valid:"-"`
	IsFollowed bool   `json:"is_followed" valid:"-"`
}
