package fake

func WithRecover(f func()) {
	defer func() {
		recover()
	}()
	f()
}
