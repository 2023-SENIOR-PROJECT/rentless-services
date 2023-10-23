package models

import "time"

type User struct {
	ID        uint      `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Email     string `json:"email" bson:"email"`
	Password1 string `json:"password1" bson:"password1"`
	Password2 string `json:"password2" bson:"password2"`

	Firstname string `json:"firstname" bson:"firstname"`
	Lastname  string `json:"lastname" bson:"lastname"`
	Age       int    `json:"age" bson:"age"`
}
type LoginRequest struct {
	Email string `json:"email" bson:"email"`
	Pwd   string `json:"pwd" bson:"pwd"`
}

type UserAuthSturct struct {
	Email   string `json:"email" bson:"email"`
	Pwd     string `json:"pwd" bson:"pwd"`
	Token   string `json:"token" bson:"token"`
	User_id uint   `json:"user_id" bson:"user_id"`
}
