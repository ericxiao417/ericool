package fake_test

import (
	"testing"

	"github.com/ericxiao417/ericool/assert"
	"github.com/ericxiao417/ericool/fake"
)

func TestWithFakeExit(t *testing.T) {
	var er fake.ExitResult
	er = fake.WithFakeOSExit(func() {
		fake.OS_Exit(1)
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)

	er = fake.WithFakeOSExit(func() {
	})
	assert.Equal(t, false, er.HasExit)

	er = fake.WithFakeOSExit(func() {
		fake.OS_Exit(2)
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 2, er.ExitCode)
}
