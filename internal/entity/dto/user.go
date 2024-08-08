package entity

type User struct {
	Id       uint32 `json:"id" valid:"-"`
	FullName string `json:"full_name" valid:"runelength(5|100)"`
	Birthday string `json:"birthday" valid:"date"`
	Email    string `json:"email" valid:"email"`
	ImgUrl   string `json:"img_url" valid:"runelength(10|100)"`
}
