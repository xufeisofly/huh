package model

import (
	"context"
	"errors"
	"time"
)

type User struct {
	Email     string
	ID        uint32    `huh:"pk"`
	CreatedAt time.Time `huh:"readonly"`
	UpdatedAt time.Time `huh:"readonly"`
}

func (u *User) TableName() string {
	return "users"
}

func (u User) BeforeCreate(ctx context.Context) error {
	if u.ID == 2 {
		return errors.New("before create error")
	}
	return nil
}

func (u *User) BeforeSave(ctx context.Context) error {
	u.Email = "update3@huh.com"
	return nil
}
