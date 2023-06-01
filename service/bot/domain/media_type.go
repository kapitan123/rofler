package domain

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

var (
	Video = MediaType{"video"}
	Image = MediaType{"image"}
)

var typeValues = []MediaType{
	Video,
	Image,
}

type MediaType struct {
	m string
}

func NewMediaTypeFromString(typeStr string) (MediaType, error) {
	val, found := lo.Find(typeValues, func(t MediaType) bool {
		return t.String() == typeStr
	})

	if found {
		return val, nil
	}

	return MediaType{}, errors.Errorf("unknown '%s' type", typeStr)
}

func (m MediaType) String() string {
	return m.m
}

func (h MediaType) IsZero() bool {
	return h == MediaType{}
}
