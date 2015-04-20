package pezauth

import "code.google.com/p/go-uuid/uuid"

//GUID interface and struct
type (
	GUIDMaker interface {
		Create() string
	}
	GUIDMake struct {
	}
)

//Create - creates a new random guid
func (s *GUIDMake) Create() string {
	r := uuid.New()
	return r
}
