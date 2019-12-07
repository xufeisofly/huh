package huh_test

import (
	"os"
	"testing"

	"github.com/xufeisofly/huh/testdata"
	model "github.com/xufeisofly/huh/testdata/models"

	"github.com/xufeisofly/huh"
)

func setup() {
	huh.Config("mysql", huh.DBConfig{
		Master: "norris@(127.0.0.1:3306)/mysite?charset=utf8",
	})
}

func tearDown() {
	huh.Close()
}

func TestEverything(t *testing.T) {
	testdata.PrepareTables()
	defer testdata.CleanUpTables()

	o := huh.New()
	ctx := huh.Context()
	user := model.User{Email: "test@huh.com", ID: 1}

	var receiver model.User
	var receivers []model.User

	// test create raw sql
	c, _ := o.Create().WithCallbacks().Of(ctx, &user)

	sqlStr := c.String()
	if sqlStr != "INSERT INTO users (email,id) VALUES ('update3@huh.com','1')" {
		t.Errorf("sqlStr actual: %s", sqlStr)
	}

	user = model.User{Email: "test@huh.com", ID: 1}
	// test normal create
	o.Create().WithCallbacks().Do(ctx, &user)

	receiver = model.User{}
	o.Get(1).Do(ctx, &receiver)
	expected := model.User{Email: "update3@huh.com", ID: 1}

	if receiver.Email != expected.Email {
		t.Errorf("[create error] expect: %v, actual: %v", expected, receiver)
	}

	// test create with hooks
	user2 := model.User{ID: 2, Email: "test2@huh.com"}
	if err := o.Create().WithCallbacks().Do(ctx, &user2); err == nil {
		t.Errorf("CreateBefore hook should have error")
	}

	receiver = model.User{}
	o.Get(2).Do(ctx, &receiver)
	if (model.User{}) != receiver {
		t.Errorf("[get error] should get nil, actual: %v", receiver)
	}

	// test update
	o.Update("email", "update@huh.com").WithCallbacks().Do(ctx, &user)

	receiver = model.User{}
	o.Get(user.ID).Do(ctx, &receiver)
	if receiver.Email != "update@huh.com" {
		t.Errorf("[update error] expected: %s, actual: %s", "update@huh.com", receiver.Email)
	}

	// test update_all with where
	user4 := model.User{ID: 4, Email: "update@huh.com"}
	o.Create().WithCallbacks().Do(ctx, &user4)

	receivers = []model.User{}
	o.Where("email = ?", "update3@huh.com").Do(ctx, &receivers)
	if len(receivers) != 1 {
		t.Errorf("[where error] result count should be 1")
	}

	receivers = []model.User{}
	o.Where("email = ?", "update@huh.com").Update("email", "update2@huh.com").Do(ctx, model.User{})
	o.Where("email = ?", "update@huh.com").Do(ctx, &receivers)

	if len(receivers) != 0 {
		t.Errorf("[where error] result count should be 0")
	}

	// test transaction no create
	user3 := model.User{ID: 3, Email: "test3@huh.com"}
	o.Transaction(ctx, func(h *huh.Orm) {
		h.Create().Do(ctx, &user3)
		h.MustCreate().Do(ctx, &user)
	})

	o.Create().WithCallbacks().Do(ctx, &user3)

	o.Where("email = ?", "update3@huh.com").Do(ctx, &receivers)

	if len(receivers) != 2 {
		t.Errorf("[where error] result count should be 2")
	}

	// test get by pk
	user5 := model.User{}
	expected = model.User{ID: uint32(1), Email: "update2@huh.com"}
	o.Get(1).Do(ctx, &user5)
	if user5.Email != expected.Email {
		t.Errorf("get error, expected: %v, actual: %v", expected, user5)
	}

	// test get by condition
	user6 := model.User{}
	expected = model.User{ID: uint32(4), Email: "update3@huh.com"}
	o.GetBy("email", "update3@huh.com", "id", 4).Do(ctx, &user6)
	if user6.Email != expected.Email {
		t.Errorf("get error, expected: %v, actual: %v", expected, user6)
	}

	user7 := model.User{}
	expected = model.User{ID: uint32(1), Email: "update2@huh.com"}
	o.GetBy(map[string]interface{}{
		"id":    1,
		"email": "update2@huh.com",
	}).Do(ctx, &user7)
	if user7.Email != expected.Email {
		t.Errorf("get error, expected: %v, actual: %v", expected, user7)
	}

	// test where
	var users = []model.User{}
	o.Where("email = ?", "update3@huh.com").Do(ctx, &users)
	expects := []model.User{
		{Email: "update3@huh.com", ID: 3},
		{Email: "update3@huh.com", ID: 4},
	}
	for i, expected := range expects {
		if users[i].Email != expected.Email {
			t.Errorf("where error, expected: %v, actual: %v", expected, users[i])
		}
	}

	// test limit offset
	users = []model.User{}
	o.Where("email = ?", "update3@huh.com").Limit(1).Offset(1).Do(ctx, &users)
	expected = model.User{Email: "update3@huh.com", ID: 4}
	if expected.Email != users[0].Email {
		t.Errorf("where error, expected: %v, actual: %v", expected, users[0])
	}

	// test order by
	users = []model.User{}
	o.Where("email = ?", "update3@huh.com").Order("id desc").Do(ctx, &users)
	expects = []model.User{
		{Email: "update3@huh.com", ID: 4},
		{Email: "update3@huh.com", ID: 3},
	}
	for i, expected := range expects {
		if users[i].Email != expected.Email {
			t.Errorf("where error, expected: %v, actual: %v", expected, users[i])
		}
	}

	// test destroy
	user = model.User{}
	o.Get(1).Do(ctx, &user)
	o.Destroy().Do(ctx, &user)

	user = model.User{}
	o.GetBy("email", "update2@huh.com").Do(ctx, &user)
	if (model.User{}) != user {
		t.Errorf("[destroy error], destroy failed")
	}

	// test destroy where
	o.Where("email = ?", "update3@huh.com").Destroy().Do(ctx, model.User{})
	users = []model.User{}
	o.Where("email = ?", "update3@huh.com").Do(ctx, &users)
	if len(users) != 0 {
		t.Errorf("[destroy error], destroy failed, expected: %v, actual: %v", []model.User{}, users)
	}

	// test nested transaction
	user8 := model.User{ID: 8, Email: "test8@huh.com"}
	user9 := model.User{ID: 9, Email: "test9@huh.com"}
	o.Transaction(ctx, func(o *huh.Orm) {
		o.Create().WithCallbacks().Do(ctx, &user8)
		// user2 cannot be created
		o.Transaction(ctx, func(o *huh.Orm) {
			o.MustCreate().WithCallbacks().Do(ctx, &user2)
		})

		o.Transaction(ctx, func(o *huh.Orm) {
			o.MustCreate().WithCallbacks().Do(ctx, &user9)
		})
	})

	users = []model.User{}
	o.Where("email = 'update3@huh.com'").Do(ctx, &users)

	if users[0].ID != 8 || users[1].ID != 9 {
		t.Errorf("[nested transaction error]")
	}

	// test multi where
	users = []model.User{}
	o.Where("email", "update3@huh.com").And("id > ?", 8).Do(ctx, &users)
	if len(users) != 1 {
		t.Errorf("[multi where error] users should be 1")
	}

	// test selected columns
	user = model.User{}
	o.Select("id", "created_at").GetBy("email", "update3@huh.com").Do(ctx, &user)
	if user.Email != "" || user.ID != 8 {
		t.Errorf("[select error] user email expected: %s, actual: %s; user id expected: %d, actual: %d", "", user.Email, 8, user.ID)
	}

	users = []model.User{}
	o.Where("id IN ?", []interface{}{8, 9}).Do(ctx, &users)
	if len(users) != 2 {
		t.Errorf("[where in error] users length should be 2")
	}

	// test where or
	users = []model.User{}
	o.Where("id", 9).Or("id", 8).Do(ctx, &users)
	if len(users) != 2 {
		t.Errorf("[where or error] users should be 2")
	}
}

func BenchmarkOf(b *testing.B) {
	var user model.User
	for i := 0; i < b.N; i++ {
		huh.New().Of(huh.Context(), &user)
	}
}

func BenchmarkSelect(b *testing.B) {
	var users []model.User
	for i := 0; i < b.N; i++ {
		huh.New().Where("id", 9).Or("id", 8).Of(huh.Context(), &users)
	}
}

func BenchmarkCreate(b *testing.B) {
	var user = model.User{}
	for i := 0; i < b.N; i++ {
		huh.New().Create().Of(huh.Context(), &user)
	}
}

func BenchmarkUpdate(b *testing.B) {
	var user = model.User{ID: 1}
	for i := 0; i < b.N; i++ {
		huh.New().Update().Of(huh.Context(), &user)
	}
}

func BenchmarkDestroy(b *testing.B) {
	var user = model.User{ID: 1}
	for i := 0; i < b.N; i++ {
		huh.New().Destroy().Of(huh.Context(), &user)
	}
}

func TestMain(m *testing.M) {
	setup()
	defer tearDown()

	os.Exit(m.Run())
}
