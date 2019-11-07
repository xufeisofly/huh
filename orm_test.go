package huh_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/xufeisofly/huh"
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

func (u *User) BeforeCreate(ctx context.Context) error {
	if u.ID == 2 {
		return errors.New("before create error")
	}
	return nil
}

func (u *User) BeforeSave(ctx context.Context) error {
	u.Email = "update3@huh.com"
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
created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY (id)
)`
	o.Exec(rawSQL)
}

func dropTable() {
	o := huh.New()
	rawSQL := `DROP TABLE users`
	o.Exec(rawSQL)
}

func TestEverything(t *testing.T) {
	dropTable()
	createTable()
	// defer dropTable()

	o := huh.New()
	ctx := huh.Context()
	user := User{Email: "test@huh.com", ID: 1}

	var receiver User
	var receivers []User

	// test create raw sql
	c, _ := o.Create().Of(ctx, &user)
	sqlStr := c.String()
	if sqlStr != "INSERT INTO users (email,id) VALUES ('update3@huh.com','1')" {
		t.Errorf("sqlStr actual: %s", sqlStr)
	}

	user = User{Email: "test@huh.com", ID: 1}
	// test normal create
	o.Create().Do(ctx, &user)

	receiver = User{}
	o.Get(1).Do(ctx, &receiver)
	expected := User{Email: "update3@huh.com", ID: 1}

	if receiver.Email != expected.Email {
		t.Errorf("[create error] expect: %v, actual: %v", expected, receiver)
	}

	// test create with hooks
	user2 := User{ID: 2, Email: "test2@huh.com"}
	if err := o.Create().Do(ctx, &user2); err == nil {
		t.Errorf("CreateBefore hook should have error")
	}

	receiver = User{}
	o.Get(2).Do(ctx, &receiver)
	if (User{}) != receiver {
		t.Errorf("[get error] should get nil, actual: %v", receiver)
	}

	// test update
	o.Update("email", "update@huh.com").Do(ctx, &user)

	receiver = User{}
	o.Get(user.ID).Do(ctx, &receiver)
	if receiver.Email != "update@huh.com" {
		t.Errorf("[update error] expected: %s, actual: %s", "update@huh.com", receiver.Email)
	}

	// test update_all with where
	user4 := User{ID: 4, Email: "update@huh.com"}
	o.Create().Do(ctx, &user4)

	receivers = []User{}
	o.Where("email = ?", "update3@huh.com").Do(ctx, &receivers)
	if len(receivers) != 1 {
		t.Errorf("[where error] result count should be 1")
	}

	receivers = []User{}
	o.Where("email = ?", "update@huh.com").Update("email", "update2@huh.com").Do(ctx, User{})
	o.Where("email = ?", "update@huh.com").Do(ctx, &receivers)

	if len(receivers) != 0 {
		t.Errorf("[where error] result count should be 0")
	}

	// test transaction no create
	user3 := User{ID: 3, Email: "test3@huh.com"}
	o.Transaction(ctx, func(h *huh.Orm) {
		h.Create().Do(ctx, &user3)
		h.MustCreate().Do(ctx, &user)
	})

	o.Create().Do(ctx, &user3)

	o.Where("email = ?", "update3@huh.com").Do(ctx, &receivers)

	if len(receivers) != 2 {
		t.Errorf("[where error] result count should be 2")
	}

	// test get by pk
	user5 := User{}
	expected = User{ID: uint32(1), Email: "update2@huh.com"}
	o.Get(1).Do(ctx, &user5)
	if user5.Email != expected.Email {
		t.Errorf("get error, expected: %v, actual: %v", expected, user5)
	}

	// test get by condition
	user6 := User{}
	expected = User{ID: uint32(4), Email: "update3@huh.com"}
	o.GetBy("email", "update3@huh.com", "id", 4).Do(ctx, &user6)
	if user6.Email != expected.Email {
		t.Errorf("get error, expected: %v, actual: %v", expected, user6)
	}

	user7 := User{}
	expected = User{ID: uint32(1), Email: "update2@huh.com"}
	o.GetBy(map[string]interface{}{
		"id":    1,
		"email": "update2@huh.com",
	}).Do(ctx, &user7)
	if user7.Email != expected.Email {
		t.Errorf("get error, expected: %v, actual: %v", expected, user7)
	}

	// test where
	var users = []User{}
	o.Where("email = ?", "update3@huh.com").Do(ctx, &users)
	expects := []User{
		{Email: "update3@huh.com", ID: 3},
		{Email: "update3@huh.com", ID: 4},
	}
	for i, expected := range expects {
		if users[i].Email != expected.Email {
			t.Errorf("where error, expected: %v, actual: %v", expected, users[i])
		}
	}

	// test limit offset
	users = []User{}
	o.Where("email = ?", "update3@huh.com").Limit(1).Offset(1).Do(ctx, &users)
	expected = User{Email: "update3@huh.com", ID: 4}
	if expected.Email != users[0].Email {
		t.Errorf("where error, expected: %v, actual: %v", expected, users[0])
	}

	// test order by
	users = []User{}
	o.Where("email = ?", "update3@huh.com").Order("id desc").Do(ctx, &users)
	expects = []User{
		{Email: "update3@huh.com", ID: 4},
		{Email: "update3@huh.com", ID: 3},
	}
	for i, expected := range expects {
		if users[i].Email != expected.Email {
			t.Errorf("where error, expected: %v, actual: %v", expected, users[i])
		}
	}

	// test destroy
	user = User{}
	o.Get(1).Do(ctx, &user)
	o.Destroy().Do(ctx, &user)

	user = User{}
	o.GetBy("email", "update2@huh.com").Do(ctx, &user)
	if (User{}) != user {
		t.Errorf("[destroy error], destroy failed")
	}

	// test destroy where
	o.Where("email = ?", "update3@huh.com").Destroy().Do(ctx, User{})
	users = []User{}
	o.Where("email = ?", "update3@huh.com").Do(ctx, &users)
	if len(users) != 0 {
		t.Errorf("[destroy error], destroy failed, expected: %v, actual: %v", []User{}, users)
	}

	// test nested transaction
	user8 := User{ID: 8, Email: "test8@huh.com"}
	user9 := User{ID: 9, Email: "test9@huh.com"}
	o.Transaction(ctx, func(o *huh.Orm) {
		o.Create().Do(ctx, &user8)
		// user2 cannot be created
		o.Transaction(ctx, func(o *huh.Orm) {
			o.MustCreate().Do(ctx, &user2)
		})

		o.Transaction(ctx, func(o *huh.Orm) {
			o.MustCreate().Do(ctx, &user9)
		})
	})

	users = []User{}
	o.Where("email = 'update3@huh.com'").Do(ctx, &users)

	if users[0].ID != 8 || users[1].ID != 9 {
		t.Errorf("[nested transaction error]")
	}

	// test multi where
	users = []User{}
	o.Where("email", "update3@huh.com").Where("id > ?", 8).Do(ctx, &users)
	if len(users) != 1 {
		t.Errorf("[multi where error] users should be 1")
	}

	// test selected columns
	user = User{}
	o.Select("id", "created_at").GetBy("email", "update3@huh.com").Do(ctx, &user)
	if user.Email != "" || user.ID != 8 {
		t.Errorf("[select error] user email expected: %s, actual: %s; user id expected: %d, actual: %d", "", user.Email, 8, user.ID)
	}

	users = []User{}
	o.Where("id IN ?", []interface{}{8, 9}).Do(ctx, &users)
	if len(users) != 2 {
		t.Errorf("[where in error] users length should be 2")
	}
}

func TestMain(m *testing.M) {
	setup()
	defer tearDown()

	os.Exit(m.Run())
}
