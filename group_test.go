package errs

import (
	"fmt"
	"strings"
	"testing"
)

func TestGroup(t *testing.T) {
	alpha := New("alpha")
	beta := New("beta")

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
		t.Fatal("expected \"group: alpha; beta\"")
	}
	if fmt.Sprintf("%v", group.Err()) != "alpha; beta" {
		t.Fatal("expected \"group: alpha; beta\"")
	}
	if strings.Count(fmt.Sprintf("%+v", group.Err()), "\n") <= 1 {
		t.Fatal("expected multiple lines with +v")
	}

	t.Logf("%%v:\n%v", group.Err())
	t.Logf("%%+v:\n%+v", group.Err())

	if Unwrap(group.Err()) != Unwrap(alpha) {
		t.Fatal("expected alpha")
	}
}
