package entity

type Subscription struct {
	FollowerId uint32
	EmployeeId uint32
	Interval   uint16
	IsFollowed bool
}
