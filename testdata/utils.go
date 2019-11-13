package testdata

import (
	"io/ioutil"
	"os"

	"github.com/xufeisofly/huh"
)

func PrepareTables() {
	o := huh.New()
	file, err := os.Open("testdata/schema.sql")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	o.Exec(string(b))
}

func CleanUpTables() {
	o := huh.New()
	rawSQL := `DROP TABLE users`
	o.Exec(rawSQL)
}
