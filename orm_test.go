package huh_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/xufeisofly/huh"
)

type User struct {
	Email string
	ID    uint32 `huh:"pk"`
}

func (u User) TableName() string {
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
	dropTable()
	createTable()
	// defer dropTable()

	o := huh.New()
	ctx := huh.Context()
	user := User{Email: "test@huh.com", ID: 1}

	// test create raw sql
	c := o.Create().Of(ctx, &user)
	sqlStr := c.String()
	if sqlStr != "INSERT INTO users (email,id) VALUES ('test@huh.com','1')" {
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

	// test update
	err = o.Update("email", "update@huh.com").Do(ctx, &user)
	if err != nil {
		t.Errorf("update error: %s", err)
	}

	// test update_all with where
	user4 := User{ID: 4, Email: "update@huh.com"}
	o.Create().Do(ctx, &user4)
	err = o.Where("email = ?", "update@huh.com").Update("email", "update2@huh.com").Do(ctx, User{})
	if err != nil {
		t.Errorf("update error: %s", err)
	}

	// test transaction
	user3 := User{ID: 3, Email: "test3@huh.com"}
	o.Transaction(ctx, func(h *huh.Orm) {
		h.Create().Do(ctx, &user3)
		h.MustCreate().Do(ctx, &user)
	})

	// // test get by pk
	// user5 := User{}
	// expected := User{ID: uint32(1), Email: "update2@huh.com"}
	// o.Get(1).Do(ctx, &user5)
	// if user5 != expected {
	// 	t.Errorf("get error, expected: %v, actual: %v", expected, user5)
	// }

	// // test get by condition
	// user6 := User{}
	// expected = User{ID: uint32(4), Email: "update2@huh.com"}
	// o.GetBy("email", "update2@huh.com", "id", 4).Do(ctx, &user6)
	// if user6 != expected {
	// 	t.Errorf("get error, expected: %v, actual: %v", expected, user6)
	// }

	// user7 := User{}
	// expected = User{ID: uint32(1), Email: "update2@huh.com"}
	// o.GetBy(map[string]interface{}{
	// 	"id":    1,
	// 	"email": "update2@huh.com",
	// }).Do(ctx, &user7)
	// if user7 != expected {
	// 	t.Errorf("get error, expected: %v, actual: %v", expected, user7)
	// }

	// test where
	var users = []User{}
	o.Where("email = ?", "update2@huh.com").Do(ctx, &users)
	expects := []User{
		{Email: "update2@huh.com", ID: 1},
		{Email: "update2@huh.com", ID: 4},
	}
	for i, expected := range expects {
		if users[i] != expected {
			t.Errorf("where error, expected: %v, actual: %v", expected, users)
		}
	}

}

func TestMain(m *testing.M) {
	setup()
	defer tearDown()

	os.Exit(m.Run())
}
