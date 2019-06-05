package errs

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/zeebo/assert"
)

func TestGroup(t *testing.T) {
	alpha := Errorf("alpha")
	beta := Errorf("beta")

	var group Group
	group.Add(nil, nil, nil)
	assert.NoError(t, group.Err())

	group.Add(alpha)
	assert.Equal(t, group.Err(), alpha)

	group.Add(nil, beta)
	assert.Equal(t, group.Err().Error(), "alpha; beta")
	assert.Equal(t, fmt.Sprintf("%v", group.Err()), "alpha; beta")
	assert.That(t, strings.Count(fmt.Sprintf("%+v", group.Err()), "\n") > 0)
	assert.Equal(t, errors.Unwrap(group.Err()), alpha)
}
