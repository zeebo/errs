// Package errs provides a simple error package with stack traces.
package errs

import (
	"fmt"
	"io"
	"runtime"
)

// Namer is implemented by all errors returned in this package. It returns a
// name for the class of error it is, and a boolean indicating if the name is
// valid.
type Namer interface{ Name() (string, bool) }

// Causer is implemented by all errors returned in this package. It returns
// the underlying cause of the error, or nil if there is no underlying cause.
type Causer interface{ Cause() error }

// New returns an error not contained in any class. This is the same as calling
// fmt.Errorf(...) except it captures a stack trace on creation.
func New(format string, args ...interface{}) error {
	return (*Class).create(nil, 3, fmt.Errorf(format, args...))
}

// Wrap returns an error not contained in any class. It just associates a stack
// trace with the error. Wrap returns nil if err is nil.
func Wrap(err error) error {
	return (*Class).create(nil, 3, err)
}

// Unwrap returns the underlying error, if any, or just the error.
func Unwrap(err error) error {
	// we call Cause as much as possible. Since comparing arbitrary interfaces
	// with equality isn't panic safe, we only loop up to 100 times to ensure
	// that a poor implementation that loops does not cause a hang.
	for i := 0; err != nil && i < 100; i++ {
		causer, ok := err.(Causer)
		if !ok {
			break
		}

		// if the cause of some error is nil, we return it.
		nerr := causer.Cause()
		if nerr == nil {
			return err
		}
		err = nerr
	}

	return err
}

// Classes returns all the classes that have wrapped the error.
func Classes(err error) []*Class {
	if err, ok := err.(*errorT); ok && err != nil {
		return append([]*Class(nil), err.classes...)
	}
	return nil
}

//
// error classes
//

// Class represents a class of errors. You can construct errors, and check if
// errors are part of the class.
type Class string

// Has returns true if the passed in error was wrapped by this class.
func (c *Class) Has(err error) bool {
	if err, ok := err.(*errorT); ok {
		for _, k := range err.classes {
			if k == c {
				return true
			}
		}
	}
	return false
}

// New constructs an error with the format string that will be contained by
// this class. This is the same as calling Wrap(fmt.Errorf(...)).
func (c *Class) New(format string, args ...interface{}) error {
	return c.create(3, fmt.Errorf(format, args...))
}

// Wrap returns a new error based on the passed in error that is contained in
// this class. Wrap returns nil if err is nil.
func (c *Class) Wrap(err error) error {
	return c.create(3, err)
}

// create constructs the error, or just adds the class to the error, keeping
// track of the stack if it needs to construct it.
func (c *Class) create(depth int, err error) error {
	if err == nil {
		return nil
	}

	if err, ok := err.(*errorT); ok {
		if c != nil && err.outerClass() != c {
			err.classes = append(err.classes, c)
		}
		return err
	}

	var pcs [256]uintptr
	n := runtime.Callers(depth, pcs[:])

	var classes []*Class
	if c != nil {
		classes = []*Class{c}
	}

	return &errorT{
		classes: classes,
		err:     err,
		pcs:     pcs[:n:n],
	}
}

//
// errors
//

// errorT is the type of errors returned from this package.
type errorT struct {
	classes []*Class
	pcs     []uintptr
	err     error
}

var ( // ensure *errorT implements the helper interfaces.
	_ Namer  = (*errorT)(nil)
	_ Causer = (*errorT)(nil)
	_ error  = (*errorT)(nil)
)

// outerClass returns the outermost wrapping class of the error.
func (e *errorT) outerClass() *Class {
	if len(e.classes) == 0 {
		return nil
	}
	return e.classes[len(e.classes)-1]
}

// errorT implements the error interface.
func (e *errorT) Error() string {
	return fmt.Sprintf("%v", e)
}

// Format handles the formatting of the error. Using a "+" on the format string
// specifier will also write the stack trace.
func (e *errorT) Format(f fmt.State, c rune) {
	for i := len(e.classes) - 1; i >= 0; i-- {
		name := string(*e.classes[i])
		if len(name) > 0 {
			fmt.Fprintf(f, "%s: ", name)
		}
	}
	fmt.Fprintf(f, "%v", e.err)

	if f.Flag(int('+')) {
		summarizeStack(f, e.pcs)
	}
}

// Cause implements the interface wrapping errors are expected to implement
// to allow getting at underlying causes.
func (e *errorT) Cause() error {
	return e.err
}

// Name returns the name for the error, which is the first wrapping class.
func (e *errorT) Name() (string, bool) {
	outer := e.outerClass()
	if outer == nil {
		return "", false
	}
	return string(*outer), true
}

// summarizeStack writes stack line entries to the writer.
func summarizeStack(w io.Writer, pcs []uintptr) {
	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		if !more {
			return
		}
		fmt.Fprintf(w, "\n\t%s:%d", frame.Function, frame.Line)
	}
}
