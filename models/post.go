package models

type Post struct {
	Post_author  string  `json:"author"`
	Post_created string  `json:"created"`
	Post_forum   string  `json:"forum"`
	Post_id      int32   `json:"id"`
	IsEdited     bool    `json:"isEdited"`
	Post_message string  `json:"message"`
	Parent       int32   `json:"parent"`
	Post_thread  int32   `json:"thread"`
	Path         []int64 `json:"-"`
	Founder      int     `json:"-"`
}

type PostInfo struct {
	Post   *Post   `json:"post"`
	Forum  *Forum  `json:"forum"`
	Thread *Thread `json:"thread"`
	Author *User   `json:"author"`
}

type PostInfoContainer struct {
	*Post
	*User
	*Thread
	*Forum
}
