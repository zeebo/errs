package errs

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
)

type causeError struct{ error }

func (c causeError) Cause() error { return c.error }

func TestErrs(t *testing.T) {
	assert := func(t *testing.T, v bool, err ...interface{}) {
		t.Helper()
		if !v {
			t.Fatal(err...)
		}
	}

	var (
		foo   = Class("foo")
		bar   = Class("bar")
		baz   = Class("baz")
		empty = Class("")
	)

	t.Run("Class", func(t *testing.T) {
		t.Run("Has", func(t *testing.T) {
			assert(t, foo.Has(foo.New("t")))
			assert(t, !foo.Has(bar.New("t")))
			assert(t, !foo.Has(baz.New("t")))

			assert(t, !bar.Has(foo.New("t")))
			assert(t, bar.Has(bar.New("t")))
			assert(t, !bar.Has(baz.New("t")))

			assert(t, foo.Has(bar.Wrap(foo.New("t"))))
			assert(t, bar.Has(bar.Wrap(foo.New("t"))))
			assert(t, !baz.Has(bar.Wrap(foo.New("t"))))

			assert(t, foo.Has(foo.Wrap(bar.New("t"))))
			assert(t, bar.Has(foo.Wrap(bar.New("t"))))
			assert(t, !baz.Has(foo.Wrap(bar.New("t"))))
		})

		t.Run("Same Name", func(t *testing.T) {
			c1 := Class("c")
			c2 := Class("c")

			assert(t, c1.Has(c1.New("t")))
			assert(t, !c2.Has(c1.New("t")))

			assert(t, !c1.Has(c2.New("t")))
			assert(t, c2.Has(c2.New("t")))
		})

		t.Run("Wrap Nil", func(t *testing.T) {
			assert(t, foo.Wrap(nil) == nil)
		})

		t.Run("WrapP", func(t *testing.T) {
			err := func() (err error) {
				defer foo.WrapP(&err)

				if 1 == 1 {
					return errors.New("err")
				}
				return nil
			}()

			t.Logf("%+v", err)
			assert(t, foo.Has(err))
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Format Contains Classes", func(t *testing.T) {
			assert(t, strings.Contains(foo.New("t").Error(), "foo"))
			assert(t, strings.Contains(bar.New("t").Error(), "bar"))

			assert(t, strings.Contains(bar.Wrap(foo.New("t")).Error(), "foo"))
			assert(t, strings.Contains(bar.Wrap(foo.New("t")).Error(), "bar"))

			assert(t, strings.Contains(foo.Wrap(bar.New("t")).Error(), "foo"))
			assert(t, strings.Contains(foo.Wrap(bar.New("t")).Error(), "bar"))
		})

		t.Run("Format With Stack", func(t *testing.T) {
			err := foo.New("t")

			assert(t,
				!strings.Contains(fmt.Sprintf("%v", err), "\n"),
				"%v format contains newline",
			)
			assert(t,
				strings.Contains(fmt.Sprintf("%+v", err), "\n"),
				"%+v format does not contain newline",
			)
		})

		t.Run("Unwrap", func(t *testing.T) {
			err := fmt.Errorf("t")

			assert(t, nil == Unwrap(nil))
			assert(t, err == Unwrap(err))
			assert(t, err == Unwrap(foo.Wrap(err)))
			assert(t, err == Unwrap(bar.Wrap(foo.Wrap(err))))
			assert(t, err == Unwrap(causeError{error: err}))

			// ensure a trivial cycle eventually completes
			loop := new(causeError)
			loop.error = loop
			assert(t, loop == Unwrap(loop))
		})

		t.Run("Cause", func(t *testing.T) {
			err := fmt.Errorf("t")

			assert(t, err == foo.Wrap(err).(*errorT).Cause())
			assert(t, err == bar.Wrap(foo.Wrap(err)).(*errorT).Cause().(*errorT).Cause())
		})

		t.Run("Classes", func(t *testing.T) {
			err := fmt.Errorf("t")
			classes := Classes(err)
			assert(t, classes == nil)

			err = foo.Wrap(err)
			classes = Classes(err)
			assert(t, len(classes) == 1)
			assert(t, classes[0] == &foo)

			err = foo.Wrap(err)
			classes = Classes(err)
			assert(t, len(classes) == 1)
			assert(t, classes[0] == &foo)

			err = bar.Wrap(err)
			classes = Classes(err)
			assert(t, len(classes) == 2)
			assert(t, classes[0] == &bar)
			assert(t, classes[1] == &foo)

			err = bar.Wrap(err)
			classes = Classes(err)
			assert(t, len(classes) == 2)
			assert(t, classes[0] == &bar)
			assert(t, classes[1] == &foo)
		})

		t.Run("Is", func(t *testing.T) {
			alpha := New("alpha")
			beta := New("beta")
			gamma := New("gamma")
			delta := New("delta")
			epsilon := New("epsilon")

			assert(t, Is(nil, nil))
			assert(t, !Is(nil, alpha))
			assert(t, Is(alpha, alpha))
			assert(t, !Is(alpha, beta))

			err := Combine(
				alpha,
				foo.Wrap(bar.Wrap(baz.Wrap(beta))),
				bar.Wrap(Combine(gamma, baz.Wrap(delta))),
			)
			assert(t, Is(err, alpha))
			assert(t, Is(err, beta))
			assert(t, Is(err, gamma))
			assert(t, Is(err, delta))
			assert(t, !Is(err, epsilon))
		})

		t.Run("IsFunc", func(t *testing.T) {
			alpha := New("alpha")
			beta := New("beta")
			gamma := New("gamma")
			delta := New("delta")
			epsilon := New("epsilon")

			assert(t, IsFunc(nil, func(err error) bool {
				return err == nil
			}))
			assert(t, !IsFunc(nil, func(err error) bool {
				return err == alpha
			}))
			assert(t, IsFunc(alpha, func(err error) bool {
				return err == alpha
			}))
			assert(t, !IsFunc(alpha, func(err error) bool {
				return err == beta
			}))

			err := Combine(
				alpha,
				foo.Wrap(bar.Wrap(baz.Wrap(beta))),
				bar.Wrap(Combine(gamma, baz.Wrap(delta))),
			)
			assert(t, IsFunc(err, func(err error) bool {
				return err == alpha
			}))
			assert(t, IsFunc(err, func(err error) bool {
				return err == beta
			}))
			assert(t, IsFunc(err, func(err error) bool {
				return err == gamma
			}))
			assert(t, IsFunc(err, func(err error) bool {
				return err == delta
			}))
			assert(t, !IsFunc(err, func(err error) bool {
				return err == epsilon
			}))
		})

		t.Run("Name", func(t *testing.T) {
			name, ok := New("t").(Namer).Name()
			assert(t, !ok)
			assert(t, name == "")

			name, ok = foo.New("t").(Namer).Name()
			assert(t, ok)
			assert(t, name == "foo")

			name, ok = bar.Wrap(foo.New("t")).(Namer).Name()
			assert(t, ok)
			assert(t, name == "bar")
		})

		t.Run("Empty String", func(t *testing.T) {
			assert(t, empty.New("test").Error() == "test")
			assert(t, foo.Wrap(empty.New("test")).Error() == "foo: test")
		})

		t.Run("Empty Format", func(t *testing.T) {
			assert(t, empty.New("").Error() == "")
			assert(t, foo.New("").Error() == "foo")
		})

		t.Run("Immutable", func(t *testing.T) {
			err := New("")
			errfoo := foo.Wrap(err)
			errbar := bar.Wrap(err)

			assert(t, err.Error() == "")
			assert(t, errfoo.Error() == "foo")
			assert(t, errbar.Error() == "bar")
		})

		t.Run("Race", func(t *testing.T) {
			err := New("race")

			var wg sync.WaitGroup
			wg.Add(2)
			go func() { foo.Wrap(err); wg.Done() }()
			go func() { bar.Wrap(err); wg.Done() }()
			wg.Wait()
		})
	})
}

func BenchmarkErrs(b *testing.B) {
	foo := Class("foo")
	err := errors.New("bench")

	b.Run("Wrap", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = foo.Wrap(err)
		}
	})

	b.Run("New", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = foo.New("bench")
		}
	})
}
