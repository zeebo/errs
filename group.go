package errs

// Group is a list of errors.
type Group []error

// Combine combines multiple errors into a single error.
func Combine(errs ...error) error {
	var group Group
	group.Add(errs...)
	return group.Err()
}

// Add adds non-empty errors to the Group.
func (group *Group) Add(errs ...error) {
	for _, err := range errs {
		if err != nil {
			*group = append(*group, err)
		}
	}
}

// Err returns an error containing all of the non-nil errors.
// If there was only one error, it will return it.
// If there were none, it returns nil.
func (group Group) Err() error {
	sanitized := group.sanitize()
	if len(sanitized) == 0 {
		return nil
	}
	if len(sanitized) == 1 {
		return sanitized[0]
	}
	return combinedError(sanitized)
}

// sanitize returns group that doesn't contain nil-s
func (group Group) sanitize() Group {
	// sanity check for non-nil errors
	for i, err := range group {
		if err == nil {
			sanitized := make(Group, 0, len(group)-1)
			sanitized = append(sanitized, group[:i]...)
			sanitized.Add(group[i+1:]...)
			return sanitized
		}
	}

	return group
}

// combinedError is a list of non-empty errors
type combinedError []error

// Cause returns the first error.
func (group combinedError) Cause() error {
	if len(group) > 0 {
		return group[0]
	}
	return nil
}

// Unwrap returns the first error.
func (group combinedError) Unwrap() error {
	return group.Cause()
}

// Error returns error string delimited by line-endings
func (group combinedError) Error() string {
	if len(group) == 0 {
		return "empty"
	}

	allErrors := group[0].Error()
	for _, err := range group[1:] {
		allErrors += "; " + err.Error()
	}
	return allErrors
}
