package fake

import "time"

var now = time.Now

func Time_Now() time.Time {
	return now()
}

func WithFakeTimeNow(n func() time.Time, fn func()) {
	now = n
	fn()
	now = time.Now
}
