package huh_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/xufeisofly/huh"
)

type User struct {
	ID    uint32 `huh:"pk"`
	Email string
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(ctx context.Context) error {
	if u.ID == 2 {
		return errors.New("before create error")
	}
	return nil
}

func setup() {
	huh.Config("mysql", huh.DBConfig{
		Master: "norris@(127.0.0.1:3306)/mysite?charset=utf8",
	})
}

func tearDown() {
	huh.Close()
}

func createTable() {
	o := huh.New()
	rawSQL := `CREATE TABLE users (
id int(11) NOT NULL AUTO_INCREMENT,
email varchar(255) NOT NULL,
PRIMARY KEY (id)
)`
	o.Exec(rawSQL)
}

func dropTable() {
	o := huh.New()
	rawSQL := `DROP TABLE users`
	o.Exec(rawSQL)
}

func TestCreate(t *testing.T) {
	createTable()
	defer dropTable()

	o := huh.New()
	ctx := huh.Context()
	user := User{Email: "test@huh.com", ID: 1}

	// test create raw sql
	c := o.Create().Of(ctx, &user)
	sqlStr := c.String()
	if sqlStr != "INSERT INTO users (id,email) VALUES ('1','test@huh.com')" {
		t.Errorf("sqlStr actual: %s", sqlStr)
	}

	// test normal create
	err := o.Create().Do(ctx, &user)
	if err != nil {
		t.Errorf("create error: %s", err)
	}

	// test create with hooks
	user2 := User{ID: 2, Email: "test2@huh.com"}
	if err := o.Create().Do(ctx, &user2); err == nil {
		t.Errorf("CreateBefore hook should have error")
	}
}

func TestMain(m *testing.M) {
	setup()
	defer tearDown()

	os.Exit(m.Run())
}
