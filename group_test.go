package errs

import (
	"errors"
	"testing"
)

func TestGroup(t *testing.T) {
	var alpha = errors.New("alpha")
	var beta = errors.New("beta")

	var group Group
	group.Add(nil, nil, nil)

	if group.Err() != nil {
		t.Fatal("expected nil")
	}

	group.Add(alpha)
	if group.Err() != alpha {
		t.Fatal("expected alpha")
	}

	group.Add(nil, beta)
	if group.Err().Error() != "alpha; beta" {
		t.Fatal("expected alpha; beta")
	}

	if Unwrap(group.Err()) != alpha {
		t.Fatal("expected alpha")
	}
}
