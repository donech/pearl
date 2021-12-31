package tabler

import (
	"reflect"
	"testing"

	"github.com/donech/pearl/pkg/db"
)

type exampleTabler struct {
	ID      int64  `columns:"ID"`
	Name    string `columns:"姓名"`
	Address string `columns:"地址"`
}

func TestColumns(t *testing.T) {
	table1 := exampleTabler{}
	table2 := struct {
		Id   string
		Name string
		Age  string
	}{}
	table3 := db.Config{}
	tests := []struct {
		name   string
		table  interface{}
		expect []string
	}{
		{"No.1", table1, []string{"ID", "姓名", "地址"}},
		{"No.2", table2, []string{"id", "name", "age"}},
		{"No.3", table3, []string{"dsn", "maxIdle", "maxOpen", "maxLifetime", "logMode"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := Columns(test.table)
			if !reflect.DeepEqual(res, test.expect) {
				t.Fatal(test.name, " failed,  got: ", res, "expect: ", test.expect)
			}
		})
	}
}
