package entity

type Notification struct {
	FollowerEmail    string `json:"follower_email"`
	FollowerId       string `json:"follower_id"`
	EmployeeEmail    string `json:"employee_email"`
	EmployeeFullName string `json:"employee_full_name"`
	EmployeeId       string `json:"employee_id"`
	Interval         uint16 `json:"interval"`
}
