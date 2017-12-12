// package errs provides a simple error package with stack traces.
package errs

import (
	"fmt"
	"io"
	"runtime"
)

// New returns an error not contained in any class. This is the same as calling
// fmt.Errorf(...) except it captures a stack trace on creation.
func New(format string, args ...interface{}) error {
	return (*Class).create(nil, 3, fmt.Errorf(format, args...))
}

// Unwrap returns the underlying error, if any, or just the error.
func Unwrap(err error) error {
	if err, ok := err.(*Error); ok && err != nil {
		return err.err
	}
	return err
}

// Classes returns all the classes that have wrapped the error.
func Classes(err error) []*Class {
	if err, ok := err.(*Error); ok && err != nil {
		return err.classes
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
	if err, ok := err.(*Error); ok {
		for _, class := range err.classes {
			if class == c {
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
// this class.
func (c *Class) Wrap(err error) error {
	return c.create(3, err)
}

// create constructs the error, or just adds the class to the error, keeping
// track of the stack if it needs to construct it.
func (c *Class) create(depth int, err error) error {
	if err, ok := err.(*Error); ok {
		if c != nil {
			err.classes = append(err.classes, c)
		}
		return err
	}

	var pcs [256]uintptr
	n := runtime.Callers(depth, pcs[:])

	var classes []*Class
	if c != nil {
		classes = append(classes, c)
	}

	return &Error{
		classes: classes,
		err:     err,
		pcs:     pcs[:n:n],
	}
}

//
// errors
//

// Error is the type of errors returned from this package.
type Error struct {
	classes []*Class
	pcs     []uintptr
	err     error
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("%v", e)
}

// Format handles the formatting of the error. Using a "+" on the format string
// specifier will also write the stack trace.
func (e *Error) Format(f fmt.State, c rune) {
	if e == nil {
		f.Write([]byte("<nil>"))
		return
	}

	for i := len(e.classes) - 1; i >= 0; i-- {
		fmt.Fprintf(f, "%s: ", string(*e.classes[i]))
	}
	fmt.Fprintf(f, "%v", e.err)

	if f.Flag(int('+')) {
		summarizeStack(f, e.pcs)
	}
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
