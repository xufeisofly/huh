package huh_test

import (
	"testing"

	"github.com/xufeisofly/huh"
	"github.com/xufeisofly/huh/testdata"
	model "github.com/xufeisofly/huh/testdata/models"
)

// 200000	      9367 ns/op
func BenchmarkSelectHandler(b *testing.B) {
	testdata.PrepareTables()
	defer testdata.CleanUpTables()

	var user model.User
	o, _ := huh.New().Of(huh.Context(), &user)

	for i := 0; i < b.N; i++ {
		huh.SelectHandler(huh.Context(), o)
	}
}

func BenchmarkCreateHandler(b *testing.B) {
	testdata.PrepareTables()
	defer testdata.CleanUpTables()

	var user = model.User{
		ID:    uint32(9),
		Email: "create@huh.com",
	}
	o, _ := huh.New().Of(huh.Context(), &user)

	for i := 0; i < b.N; i++ {
		huh.CreateHandler(huh.Context(), o)
	}
}

func BenchmarkBeforeCreateHandler(b *testing.B) {
	testdata.PrepareTables()
	defer testdata.CleanUpTables()

	var user = model.User{
		ID:    uint32(9),
		Email: "create@huh.com",
	}
	o, _ := huh.New().Of(huh.Context(), &user)

	for i := 0; i < b.N; i++ {
		huh.BeforeCreateHandler(huh.Context(), o)
	}
}
