package models

import (
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID   uuid.UUID `json:"user_id" db:"user_id" validate:"omitempty"`
	Login    string    `json:"login" db:"login" validate:"required,lte=50"`
	Password string    `json:"password,omitempty" db:"password" validate:"omitempty,required"`
	Role     string    `json:"role,omitempty" db:"role" validate:"omitempty,lte=20"`
}

type UserList struct {
	TotalCount int    `json:"total_count"`
	TotalPages int    `json:"total_pages"`
	Page       int    `json:"page"`
	Size       int    `json:"size"`
	HasMore    bool   `json:"has_more"`
	Users      []User `json:"users"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (u *User) SanitizePassword() {
	u.Password = ""
}

func (u *User) PrepareCreate() error {
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}

	if u.Role != "" {
		u.Role = strings.ToLower(strings.TrimSpace(u.Role))
	}
	return nil
}

func (u *User) PrepareUpdate() error {

	if u.Role != "" {
		u.Role = strings.ToLower(strings.TrimSpace(u.Role))
	}
	return nil
}

type UserWithToken struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
