package fake_test

import (
	"testing"
	"time"

	"github.com/ericxiao417/ericool/assert"

	"github.com/ericxiao417/ericool/fake"
)

func TestWithFakeTimeNow(t *testing.T) {
	fake.WithFakeTimeNow(func() time.Time {
		return time.Now().Add(time.Duration(2 * time.Hour))
	}, func() {
		n := fake.Time_Now()
		assert.Equal(t, true, n.Sub(time.Now()).Hours() > 1)
	})
}
