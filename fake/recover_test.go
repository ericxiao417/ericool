package fake_test

import (
	"testing"

	"github.com/ericxiao417/ericool/fake"
)

func TestWithRecover(t *testing.T) {
	fake.WithRecover(func() {
		// noop
	})
	fake.WithRecover(func() {
		panic(0)
	})
}
