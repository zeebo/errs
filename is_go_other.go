//go:build !go1.20

package errs

// Is checks if any of the underlying errors matches target
func Is(err, target error) bool {
	return IsFunc(err, func(err error) bool { return err == target })
}
