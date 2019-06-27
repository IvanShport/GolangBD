package models

type Vote struct {
	Nickname string `json:"nickname"`
	Thread   int32  `json:"thread"`
	Voice    int32  `json:"voice"`
}
