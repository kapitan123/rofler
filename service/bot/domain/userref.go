package domain

import "fmt"

type UserRef struct {
	DisplayName string
	Id          int64
}

func NewUserRef(id int64, firstName string, lastName string) UserRef {
	return UserRef{
		Id:          id,
		DisplayName: fmt.Sprintf("%s %s", firstName, lastName),
	}
}

func (ur UserRef) AsUserMention() string {
	return fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>", ur.Id, ur.DisplayName)
}
