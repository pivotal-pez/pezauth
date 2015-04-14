package pezauth

import "code.google.com/p/go-uuid/uuid"

type (
	GUIDMaker interface {
		Create() string
	}
	GUIDMake struct {
	}
)

func (s *GUIDMake) Create() string {
	r := uuid.NewRandom()
	return string(r[:])
}
