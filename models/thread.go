package models

import "time"

type Thread struct {
	Thread_author  string     `json:"author"`
	Thread_created *time.Time `json:"created"`
	Thread_forum   string     `json:"forum"`
	Thread_id      int32      `json:"id"`
	Thread_message string     `json:"message"`
	Thread_slug    *string    `json:"slug"`
	Thread_title   string     `json:"title"`
	Votes          int32      `json:"votes"`
}
