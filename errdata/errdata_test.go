package errdata

import (
	"testing"

	"github.com/zeebo/errs"
)

func TestErrdata(t *testing.T) {
	assert := func(t *testing.T, v bool, err ...interface{}) {
		t.Helper()
		if !v {
			t.Fatal(err...)
		}
	}

	var (
		foo = errs.Class("foo")
		bar = errs.Class("bar")
		baz = errs.Class("baz")
	)

	type key1 struct{}
	type key2 struct{}

	Set(&foo, key1{}, "foo 1")
	Set(&foo, key2{}, "foo 2")
	Set(&bar, key1{}, "bar 1")
	Set(&bar, key2{}, "bar 2")
	Set(&baz, key1{}, "baz 1")
	Set(&baz, key2{}, "baz 2")

	assert(t, Get(errs.New("t"), key1{}) == nil)
	assert(t, Get(errs.New("t"), key2{}) == nil)

	assert(t, Get(foo.New("t"), key1{}) == "foo 1")
	assert(t, Get(foo.New("t"), key2{}) == "foo 2")

	assert(t, Get(bar.New("t"), key1{}) == "bar 1")
	assert(t, Get(bar.New("t"), key2{}) == "bar 2")

	assert(t, Get(baz.New("t"), key1{}) == "baz 1")
	assert(t, Get(baz.New("t"), key2{}) == "baz 2")

	assert(t, Get(foo.Wrap(baz.New("t")), key1{}) == "foo 1")
	assert(t, Get(foo.Wrap(baz.New("t")), key2{}) == "foo 2")

	assert(t, Get(bar.Wrap(foo.Wrap(baz.New("t"))), key1{}) == "bar 1")
	assert(t, Get(bar.Wrap(foo.Wrap(baz.New("t"))), key2{}) == "bar 2")

	Set(&foo, key1{}, nil)
	Set(&foo, key2{}, nil)

	assert(t, Get(foo.New("t"), key1{}) == nil)
	assert(t, Get(foo.New("t"), key2{}) == nil)
}
