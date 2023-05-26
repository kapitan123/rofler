package domain

import "fmt"

type UserRef struct {
	DisplayName string
	Id          int64
}

func (ur UserRef) AsUserMention() string {
	return fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>", ur.Id, ur.DisplayName)
}
