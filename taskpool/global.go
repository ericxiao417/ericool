package taskpool

var global Pool

func Go(task TaskFn, param ...interface{}) {
	global.Go(task, param)
}

func GetCurrentStatus() Status {
	return global.GetCurrentStatus()
}

func KillIdleWorkers() {
	global.KillIdleWorkers()
}

func Init(modOptions ...ModOption) error {
	var err error
	global, err = NewPool(modOptions...)
	return err
}

func init() {
	_ = Init()
}
