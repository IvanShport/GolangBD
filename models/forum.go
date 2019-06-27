package models

type Forum struct {
	Posts       int32  `json:"posts"`
	Forum_slug  string `json:"slug"`
	Threads     int32  `json:"threads"`
	Forum_title string `json:"title"`
	Forum_user  string `json:"user"`
}
