package huh_test

import (
	"context"
	"os"
	"testing"

	"github.com/xufeisofly/huh"
)

type User struct {
	ID    uint32
	Email string
}

func (u *User) TableName() string {
	return "users"
}

func setup() {
	huh.Config("mysql", huh.DBConfig{
		Master: "root@127.0.0.1/mysite?charset=utf8&parseTime=True&loc=local",
	})
}

func tearDown() {
	huh.Close()
}

func TestCreateOf(t *testing.T) {
	o := huh.New()
	user := User{ID: 1, Email: "test@huh.com"}

	c := o.Create().Of(context.TODO(), &user)

	sqlStr := c.String()
	if sqlStr != "INSERT INTO users (ID,Email) VALUES (1,test@huh.com)" {
		t.Errorf("sqlStr actual: %s", sqlStr)
	}
}

func TestMain(m *testing.M) {
	setup()
	defer tearDown()
	os.Exit(m.Run())
}
