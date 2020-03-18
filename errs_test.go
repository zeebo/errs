package errs

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/zeebo/assert"
)

func TestErrs(t *testing.T) {
	const (
		foo   = Tag("foo")
		bar   = Tag("bar")
		baz   = Tag("baz")
		empty = Tag("")
	)

	t.Run("Tag", func(t *testing.T) {
		t.Run("Is", func(t *testing.T) {
			assert.That(t, errors.Is(foo.Errorf("t"), foo))
			assert.That(t, !errors.Is(bar.Errorf("t"), foo))
			assert.That(t, !errors.Is(baz.Errorf("t"), foo))

			assert.That(t, !errors.Is(foo.Errorf("t"), bar))
			assert.That(t, errors.Is(bar.Errorf("t"), bar))
			assert.That(t, !errors.Is(baz.Errorf("t"), bar))

			assert.That(t, errors.Is(bar.Wrap(foo.Errorf("t")), foo))
			assert.That(t, errors.Is(bar.Wrap(foo.Errorf("t")), bar))
			assert.That(t, !errors.Is(bar.Wrap(foo.Errorf("t")), baz))

			assert.That(t, errors.Is(foo.Wrap(bar.Errorf("t")), foo))
			assert.That(t, errors.Is(foo.Wrap(bar.Errorf("t")), bar))
			assert.That(t, !errors.Is(foo.Wrap(bar.Errorf("t")), baz))
		})

		t.Run("Same Name", func(t *testing.T) {
			t1 := Tag("c")
			t2 := Tag("c")

			assert.That(t, errors.Is(t1.Errorf("t"), t1))
			assert.That(t, errors.Is(t1.Errorf("t"), t2))
			assert.That(t, errors.Is(t2.Errorf("t"), t1))
			assert.That(t, errors.Is(t2.Errorf("t"), t2))
		})

		t.Run("Wrap", func(t *testing.T) {
			assert.That(t, foo.Wrap(nil) == nil)
		})

		t.Run("Empty", func(t *testing.T) {
			assert.Equal(t, empty.Errorf("test").Error(), `test`)
			assert.Equal(t, foo.Wrap(empty.Errorf("test")).Error(), `foo: test`)
			assert.Equal(t, empty.Errorf("").Error(), "")
			assert.Equal(t, foo.Errorf("").Error(), "foo")
			assert.That(t, errors.Is(empty.Wrap(foo), empty))
		})

		t.Run("Hoist", func(t *testing.T) {
			assert.Equal(t,
				foo.Errorf("context: %w", foo.Wrap(Errorf("test"))).Error(),
				`foo: context: test`)
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Format", func(t *testing.T) {
			assert.That(t, strings.Contains(foo.Errorf("t").Error(), "foo"))
			assert.That(t, strings.Contains(bar.Errorf("t").Error(), "bar"))

			assert.That(t, strings.Contains(bar.Wrap(foo.Errorf("t")).Error(), "foo"))
			assert.That(t, strings.Contains(bar.Wrap(foo.Errorf("t")).Error(), "bar"))

			assert.That(t, strings.Contains(foo.Wrap(bar.Errorf("t")).Error(), "foo"))
			assert.That(t, strings.Contains(foo.Wrap(bar.Errorf("t")).Error(), "bar"))
		})

		t.Run("Format With Stack", func(t *testing.T) {
			err := foo.Errorf("t")

			assert.That(t, !strings.Contains(fmt.Sprintf("%v", err), "\n"))
			assert.That(t, strings.Contains(fmt.Sprintf("%+v", err), "\n"))
		})

		t.Run("Unwrap", func(t *testing.T) {
			err := fmt.Errorf("t")

			assert.Equal(t, nil, errors.Unwrap(nil))
			assert.Equal(t, err, errors.Unwrap(foo.Wrap(err)))
			assert.Equal(t, err, errors.Unwrap(bar.Wrap(err)))
		})

		t.Run("Tags", func(t *testing.T) {
			err := fmt.Errorf("t")
			tags := Tags(err)
			assert.That(t, tags == nil)

			err = foo.Wrap(err)
			tags = Tags(err)
			assert.Equal(t, len(tags), 1)
			assert.Equal(t, tags[0], foo)

			err = foo.Wrap(err)
			tags = Tags(err)
			assert.Equal(t, len(tags), 1)
			assert.Equal(t, tags[0], foo)

			err = bar.Wrap(err)
			tags = Tags(err)
			assert.Equal(t, len(tags), 2)
			assert.Equal(t, tags[0], bar)
			assert.Equal(t, tags[1], foo)

			err = bar.Wrap(err)
			tags = Tags(err)
			assert.Equal(t, len(tags), 2)
			assert.Equal(t, tags[0], bar)
			assert.Equal(t, tags[1], foo)
		})

		t.Run("Name", func(t *testing.T) {
			type Namer interface{ Name() (string, bool) }

			name, ok := Errorf("t").(Namer).Name()
			assert.That(t, !ok)
			assert.Equal(t, name, "")

			name, ok = foo.Errorf("t").(Namer).Name()
			assert.That(t, ok)
			assert.Equal(t, name, "foo")

			name, ok = bar.Wrap(foo.Errorf("t")).(Namer).Name()
			assert.That(t, ok)
			assert.Equal(t, name, "bar")
		})

		t.Run("Immutable", func(t *testing.T) {
			err := Errorf("")
			errfoo := foo.Wrap(err)
			errbar := bar.Wrap(err)

			assert.Equal(t, err.Error(), "")
			assert.Equal(t, errfoo.Error(), "foo")
			assert.Equal(t, errbar.Error(), "bar")
		})

		t.Run("Race", func(t *testing.T) {
			err := Errorf("race")

			var wg sync.WaitGroup
			wg.Add(2)
			go func() { _ = foo.Wrap(err); wg.Done() }()
			go func() { _ = bar.Wrap(err); wg.Done() }()
			wg.Wait()
		})
	})
}

func BenchmarkErrs(b *testing.B) {
	foo := Tag("foo")
	err := errors.New("bench")

	b.Run("Wrap", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = foo.Wrap(err)
		}
	})

	b.Run("Errorf", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = foo.Errorf("bench")
		}
	})
}
