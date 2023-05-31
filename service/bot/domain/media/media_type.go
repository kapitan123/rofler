package media

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

var (
	Video = Type{"video"}
	Image = Type{"image"}
)

var typeValues = []Type{
	Video,
	Image,
}

type Type struct {
	m string
}

func NewTypeFromString(typeStr string) (Type, error) {
	val, found := lo.Find(typeValues, func(t Type) bool {
		return t.String() == typeStr
	})

	if found {
		return val, nil
	}

	return Type{}, errors.Errorf("unknown '%s' type", typeStr)
}

func (m Type) String() string {
	return m.m
}

func (h Type) IsZero() bool {
	return h == Type{}
}
