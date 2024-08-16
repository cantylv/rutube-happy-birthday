package entity

type SignUpForm struct {
	FullName string `json:"full_name" valid:"runelength(5|100)"`
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"password" bson:"password"`
	Birthday string `json:"birthday" valid:"date"`
}

type SignInForm struct {
	Email    string `json:"email" valid:"email" bson:"email"`
	Password string `json:"password" valid:"password" bson:"password"`
}
