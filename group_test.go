package errs

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/zeebo/assert"
)

func TestGroup(t *testing.T) {
	const (
		foo = Tag("foo")
		bar = Tag("bar")
		baz = Tag("baz")
	)

	alpha := Errorf("alpha")
	beta := Errorf("beta")
	gamma := Errorf("gamma")
	delta := Errorf("delta")
	epsilon := Errorf("epsilon")

	t.Run("Empty", func(t *testing.T) {
		var group Group
		group.Add(nil, nil, nil)
		assert.NoError(t, group.Err())
	})

	t.Run("Single", func(t *testing.T) {
		var group Group
		group.Add(alpha, nil)
		assert.Equal(t, group.Err(), alpha)
	})

	t.Run("Multiple", func(t *testing.T) {
		var group Group
		group.Append(alpha)
		group.Add(nil, beta)
		err := group.Err()

		assert.Equal(t, err.Error(), "alpha; beta")
		assert.Equal(t, fmt.Sprintf("%v", err), "alpha; beta")
		assert.That(t, strings.Count(fmt.Sprintf("%+v", err), "\n") > 0)
		assert.Equal(t, errors.Unwrap(err), alpha)
	})

	t.Run("Name", func(t *testing.T) {
		type Namer interface{ Name() (string, bool) }

		err := Combine(Tagged("t2", beta), Tagged("t1", alpha))
		assert.Equal(t, err.Error(), "t2: beta; t1: alpha")

		name, ok := err.(Namer).Name()
		assert.That(t, ok)
		assert.Equal(t, name, "group: t1; t2")
	})

	t.Run("Is", func(t *testing.T) {
		err := Combine(
			alpha,
			foo.Wrap(bar.Wrap(baz.Wrap(beta))),
			bar.Wrap(Combine(gamma, baz.Wrap(delta))),
		)

		assert.That(t, errors.Is(err, alpha))
		assert.That(t, errors.Is(err, beta))
		assert.That(t, errors.Is(err, gamma))
		assert.That(t, errors.Is(err, delta))
		assert.That(t, !errors.Is(err, epsilon))
	})
}
