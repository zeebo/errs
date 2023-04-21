package errs

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Group is a list of errors.
type Group []error

// Combine combines multiple non-empty errors into a single error.
func Combine(errs ...error) error {
	var group Group
	group.Add(errs...)
	return group.Err()
}

// Add adds non-empty errors to the Group.
func (g *Group) Add(errs ...error) {
	for _, err := range errs {
		g.Append(err)
	}
}

// Append adds a non-empty error to the Group.
func (g *Group) Append(err error) {
	if err != nil {
		*g = append(*g, err)
	}
}

// Err returns an error containing all of the non-nil errors.
// If there was only one error, it will return it.
// If there were none, it returns nil.
func (g Group) Err() error {
	sanitized := g.sanitize()
	if len(sanitized) == 0 {
		return nil
	}
	if len(sanitized) == 1 {
		return sanitized[0]
	}
	return groupedErrors(sanitized)
}

// sanitize returns group that doesn't contain nil-s
func (g Group) sanitize() Group {
	// sanity check for non-nil errors
	for i, err := range g {
		if err == nil {
			sanitized := make(Group, 0, len(g)-1)
			sanitized = append(sanitized, g[:i]...)
			sanitized.Add(g[i+1:]...)
			return sanitized
		}
	}
	return g
}

// groupedErrors is a list of non-empty errors
type groupedErrors []error

// Cause returns the first error.
func (g groupedErrors) Cause() error {
	if len(g) > 0 {
		return g[0]
	}
	return nil
}

// Unwrap returns the all of the errors.
func (g groupedErrors) Unwrap() []error { return g }

// Ungroup returns all errors.
func (g groupedErrors) Ungroup() []error { return g }

// Is is for go1.13 errors so that the Is function reports true if the error is
// part of the group.
func (g groupedErrors) Is(target error) bool {
	for _, err := range g {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// Error returns error string delimited by semicolons.
func (g groupedErrors) Error() string { return fmt.Sprintf("%v", g) }

// Name returns the set of names in the group in sorted order so that it is
// stable.
func (g groupedErrors) Name() (string, bool) {
	var names []string
	for _, err := range g {
		if namer, ok := err.(interface{ Name() (string, bool) }); ok {
			if name, ok := namer.Name(); ok {
				names = append(names, name)
			}
		}
	}
	if len(names) == 0 {
		return "group", true
	}
	sort.Strings(names)
	return "group: " + strings.Join(names, "; "), true
}

// Format handles the formatting of the error. Using a "+" on the format
// string specifier will cause the errors to be formatted with "+" and
// delimited by newlines. They are delimited by semicolons otherwise.
func (g groupedErrors) Format(f fmt.State, c rune) {
	delim := "; "
	if f.Flag(int('+')) {
		_, _ = io.WriteString(f, "group:\n--- ")
		delim = "\n--- "
	}

	for i, err := range g {
		if i != 0 {
			_, _ = io.WriteString(f, delim)
		}
		if formatter, ok := err.(fmt.Formatter); ok {
			formatter.Format(f, c)
		} else {
			fmt.Fprintf(f, "%v", err)
		}
	}
}
