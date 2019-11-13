package huh_test

import (
	"testing"

	"github.com/xufeisofly/huh"
	model "github.com/xufeisofly/huh/testdata/models"
)

func BenchmarkProcess(b *testing.B) {
	var user = model.User{
		ID:    uint32(9),
		Email: "create@huh.com",
	}
	o, _ := huh.New().Of(huh.Context(), &user)
	cb := huh.CreateCallback

	for i := 0; i < b.N; i++ {
		cb.Processor.Process(huh.Context(), o)
	}
}

func BenchmarkCallCallbacks(b *testing.B) {
	var user = model.User{
		ID:    uint32(9),
		Email: "create@huh.com",
	}
	o, _ := huh.New().Create().Of(huh.Context(), &user)

	for i := 0; i < b.N; i++ {
		o.CallCallbacks(huh.Context())
	}
}
