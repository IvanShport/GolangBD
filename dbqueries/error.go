package dbqueries

import (
	"errors"
	"fmt"
)

var ErrParentPost = errors.New("Parent post was created in another thread")

type NotFound struct {
	Model  string
	Params string
}

func (s NotFound) Error() string {
	return fmt.Sprintf(`%s error: record with "%s" not found`, s.Model, s.Params)
}
