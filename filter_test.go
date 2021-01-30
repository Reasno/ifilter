package ifilter

import (
	"errors"
	"reflect"
	"testing"
)

type mock struct {}

func (mock) f() {}

type i interface {
	f()
}

type mock2 struct {}

func (mock2) f() {}
func (mock2) f2() {}

func TestCollection_Filter(t *testing.T) {
	cases := []struct{
		name string
		data []interface{}
		target interface{}
		shouldErr bool
	}{
		{
			"basic",
			[]interface{}{mock{}, struct{}{}},
			func(filtered []i) {
				if len(filtered) != 1 {
					t.Fatal()
				}
				if ! reflect.DeepEqual(filtered[0],mock{}) {
					t.Fatal("mock is not the same")
				}
			},
			false,
		},
		{
			"none",
			[]interface{}{struct{}{}},
			func(filtered []i) {
				if len(filtered) != 0 {
					t.Fatal()
				}
			},
			false,
		},
		{
			"all",
			[]interface{}{mock{}, mock{}},
			func(filtered []i) {
				if len(filtered) != 2 {
					t.Fatal()
				}
			},
			false,
		},
		{
			"empty",
			[]interface{}{},
			func(filtered []i) {
				if len(filtered) != 0 {
					t.Fatal()
				}
			},
			false,
		},
		{
			"two mock of different type",
			[]interface{}{mock{}, mock2{}},
			func(filtered []i) {
				if len(filtered) != 2 {
					t.Fatal()
				}
			},
			false,
		},
		{
			"return err from callback",
			[]interface{}{},
			func(filtered []i) error {
				return errors.New("err")
			},
			true,
		},
		{
			"not valid interface",
			[]interface{}{mock{}, struct{}{}},
			func(filtered []mock) {},
			true,
		},
		{
			"not valid interface",
			[]interface{}{mock{}, struct{}{}},
			func(filtered int) {},
			true,
		},
		{
			"not valid interface",
			[]interface{}{mock{}, struct{}{}},
			func() {},
			true,
		},
		{
			"not valid interface",
			[]interface{}{mock{}, struct{}{}},
			1,
			true,
		},
		{
			"not valid interface",
			[]interface{}{mock{}, struct{}{}},
			nil,
			true,
		},
	}

	for _, cc := range cases {
		c := cc
		t.Run(c.name, func(t *testing.T) {
			col := Collection(c.data)
			err := col.Filter(c.target)
			if c.shouldErr && err == nil {
				t.Fatal("should err, but got none")
			}
			if ! c.shouldErr && err != nil {
				t.Fatalf("should not err, but got %s", err)
			}
		})
	}
}
