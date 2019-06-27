package models

type User struct {
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	About    string `json:"about"`
}
