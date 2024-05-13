package models

import (
	"time"
)

// User defines user model for all kinds of user in this system
type User struct {
	DisplayUser
	Type      UserType  `json:"-" db:"type"`
	PushToken *string   `json:"-" db:"push_token"`
	Email     *string   `json:"email,omitempty" db:"email"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

type UserType int

const (
	Audience UserType = iota
	Official
	Admin
)

func (ut UserType) String() string {
	switch ut {
	case Audience:
		return "audience"
	case Official:
		return "official"
	case Admin:
		return "admin"
	}
	return ""
}

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

func (g *Gender) String() string {
	return string(*g)
}

type UserRole string

const (
	Student    UserRole = "student"
	Doctor     UserRole = "doctor"
	Nurse      UserRole = "nurse"
	Pharmacist UserRole = "pharmacist"
)

func (r *UserRole) String() string {
	return string(*r)
}

// DisplayUser is a much simpler struct for display purpose
type DisplayUser struct {
	ID      string  `json:"user_id" db:"id" example:"8dfd0f04-c379-4a18-ac1b-b5c28c70d9e3"`
	Name    string  `json:"name" db:"name" example:"York Chou"`
	Picture *string `json:"picture" db:"picture"`
}
