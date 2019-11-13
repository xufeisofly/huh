package huh_test

import (
	"testing"
	"time"

	"github.com/xufeisofly/huh"
	model "github.com/xufeisofly/huh/testdata/models"
)

func TestGetModel(t *testing.T) {
	var in model.User
	var model *huh.Model

	model = huh.GetModel(&in)

	if model.TableName != "users" {
		t.Errorf("[GetModel] TableName expected: %s, actual: %s", "users", model.TableName)
	}
	if model.PrimaryField.Name != "ID" {
		t.Errorf("[GetModel] Primary Key expected: %s, actual: %s", "ID", model.PrimaryField.Name)
	}

	columnFields := model.Columns()
	expectedFields := []string{"email", "id", "created_at", "updated_at"}
	for i := range columnFields {
		if columnFields[i] != expectedFields[i] {
			t.Errorf("[GetModel] field not equal, expected: %s, actual: %s", expectedFields[i], columnFields[i])
		}
	}

	columnFields = model.WritableColumns()
	expectedFields = []string{"email", "id"}
	for i := range columnFields {
		if columnFields[i] != expectedFields[i] {
			t.Errorf("[GetModel] field not equal, expected: %s, actual: %s", expectedFields[i], columnFields[i])
		}
	}

	var initTime time.Time
	columnValues := model.Values()
	expectedValues := []interface{}{"", uint32(0), initTime, initTime}
	for i := range columnValues {
		if columnValues[i] != expectedValues[i] {
			t.Errorf("[GetModel] field %v not equal, expected: %v, actual: %v", i, expectedValues[i], columnValues[i])
		}
	}
}

// 300000	      5480 ns/op
func BenchmarkGetModel(b *testing.B) {
	var in model.User
	for i := 0; i < b.N; i++ {
		huh.GetModel(&in)
	}
}
