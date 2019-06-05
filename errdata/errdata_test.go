package errdata

import (
	"testing"

	"github.com/zeebo/assert"
	"github.com/zeebo/errs/v2"
)

func TestErrdata(t *testing.T) {
	var (
		foo = errs.Tag("foo")
		bar = errs.Tag("bar")
		baz = errs.Tag("baz")
	)

	type key1 struct{}
	type key2 struct{}

	Set(foo, key1{}, "foo 1")
	Set(foo, key2{}, "foo 2")
	Set(bar, key1{}, "bar 1")
	Set(bar, key2{}, "bar 2")
	Set(baz, key1{}, "baz 1")
	Set(baz, key2{}, "baz 2")

	assert.Equal(t, Get(errs.Errorf("t"), key1{}), nil)
	assert.Equal(t, Get(errs.Errorf("t"), key2{}), nil)

	assert.Equal(t, Get(foo.Errorf("t"), key1{}), "foo 1")
	assert.Equal(t, Get(foo.Errorf("t"), key2{}), "foo 2")

	assert.Equal(t, Get(bar.Errorf("t"), key1{}), "bar 1")
	assert.Equal(t, Get(bar.Errorf("t"), key2{}), "bar 2")

	assert.Equal(t, Get(baz.Errorf("t"), key1{}), "baz 1")
	assert.Equal(t, Get(baz.Errorf("t"), key2{}), "baz 2")

	assert.Equal(t, Get(foo.Wrap(baz.Errorf("t")), key1{}), "foo 1")
	assert.Equal(t, Get(foo.Wrap(baz.Errorf("t")), key2{}), "foo 2")

	assert.Equal(t, Get(bar.Wrap(foo.Wrap(baz.Errorf("t"))), key1{}), "bar 1")
	assert.Equal(t, Get(bar.Wrap(foo.Wrap(baz.Errorf("t"))), key2{}), "bar 2")

	Set(foo, key1{}, nil)
	Set(foo, key2{}, nil)

	assert.Equal(t, Get(foo.Errorf("t"), key1{}), nil)
	assert.Equal(t, Get(foo.Errorf("t"), key2{}), nil)
}
