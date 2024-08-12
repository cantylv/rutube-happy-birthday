package entity

type User struct {
	Id       string         `json:"id" valid:"-" bson:"_id"`
	FullName string         `json:"full_name" valid:"runelength(5|100)" bson:"full_name"`
	Birthday string         `json:"birthday" valid:"date" bson:"birthday"`
	Email    string         `json:"email" valid:"email" bson:"email"`
	ImgUrl   string         `json:"img_url" valid:"runelength(10|100)" bson:"img_url"`
	Subs     []Subscription `json:"subs" valid:"-" bson:"subs,omitempty"`
}
