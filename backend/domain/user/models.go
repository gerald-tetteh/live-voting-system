package user

import "time"

type UserRole string

const (
	Admin UserRole = "admin"
	Base UserRole = "base"
)

type User struct {
	ID string
	FirstName string
	LastName string
	MiddleName string
	Email string
	Role UserRole
	Active bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (role UserRole) IsValid() bool {
	switch role {
	case Admin, Base:
		return true
	}
	return false
}