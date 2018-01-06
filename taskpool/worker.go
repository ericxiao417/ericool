package taskpool

type worker struct {
	taskChan chan taskWrapper
	p        *pool
}

func NewWorker(p *pool) *worker {
	return &worker{
		taskChan: make(chan taskWrapper, 1),
		p:        p,
	}
}

func (w *worker) Start() {
	go func() {
		for {
			task, ok := <-w.taskChan
			if !ok {
				break
			}
			task.taskFn(task.param...)
			w.p.onIdle(w)
		}
	}()
}

func (w *worker) Stop() {
	close(w.taskChan)
}

func (w *worker) Go(t taskWrapper) {
	w.taskChan <- t
}
