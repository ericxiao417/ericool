package taskpool_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ericxiao417/ericool/taskpool"

	"github.com/ericxiao417/ericool/assert"
	"github.com/ericxiao417/ericool/ericlog"
)

var (
	taskNum       = 1000 * 1000
	initWorkerNum = 1 //1000 * 20 //1000 * 10
)

func BenchmarkOriginGo(b *testing.B) {
	ericlog.Debug("> BenchmarkOriginGo")
	var wg sync.WaitGroup
	for j := 0; j < 1; j++ {
		wg.Add(taskNum)
		for i := 0; i < taskNum; i++ {
			go func() {
				time.Sleep(10 * time.Millisecond)
				wg.Done()
			}()
		}
		wg.Wait()
	}
	ericlog.Debug("< BenchmarkOriginGo")
}

func BenchmarkTaskPool(b *testing.B) {
	ericlog.Debug("> BenchmarkTaskPool")
	var wg sync.WaitGroup
	p, _ := taskpool.NewPool(func(option *taskpool.Option) {
		option.InitWorkerNum = initWorkerNum
	})

	b.ResetTimer()
	for j := 0; j < 1; j++ {
		//b.StartTimer()
		wg.Add(taskNum)
		for i := 0; i < taskNum; i++ {
			p.Go(func(param ...interface{}) {
				time.Sleep(10 * time.Millisecond)
				wg.Done()
			})
		}
		wg.Wait()
	}
	ericlog.Debug("< BenchmarkTaskPool")
}

func TestTaskPool(t *testing.T) {
	var wg sync.WaitGroup
	p, _ := taskpool.NewPool(func(option *taskpool.Option) {
		option.InitWorkerNum = 1
	})

	go func() {
		//for {
		ericlog.Debugf("timer, worker num. status=%+v", p.GetCurrentStatus())
		time.Sleep(10 * time.Millisecond)
		//}
	}()

	n := 1000
	wg.Add(n)
	ericlog.Debug("start.")
	for i := 0; i < n; i++ {
		p.Go(func(param ...interface{}) {
			time.Sleep(10 * time.Millisecond)
			wg.Done()
		})
	}
	wg.Wait()
	ericlog.Debugf("done, worker num. status=%+v", p.GetCurrentStatus()) // 此时还有个别busy也是正常的，因为只是业务方的任务代码执行完了，可能还没回收到idle队列中
	p.KillIdleWorkers()
	ericlog.Debugf("killed, worker num. status=%+v", p.GetCurrentStatus())

	time.Sleep(100 * time.Millisecond)

	wg.Add(n)
	for i := 0; i < n; i++ {
		p.Go(func(param ...interface{}) {
			time.Sleep(10 * time.Millisecond)
			wg.Done()
		})
	}
	wg.Wait()
	ericlog.Debugf("done, worker num. status=%+v", p.GetCurrentStatus())
}

func TestMaxWorker(t *testing.T) {
	p, err := taskpool.NewPool(func(option *taskpool.Option) {
		option.MaxWorkerNum = 128
	})
	assert.Equal(t, nil, err)

	go func() {
		for i := 0; i < 5; i++ {
			ericlog.Debugf("timer. status=%+v", p.GetCurrentStatus())
			time.Sleep(100 * time.Millisecond)
		}
	}()

	var wg sync.WaitGroup
	var sum int32
	n := 1000
	wg.Add(n)
	ericlog.Debugf("start.")
	for i := 0; i < n; i++ {
		p.Go(func(param ...interface{}) {
			a := param[0].(int)
			b := param[1].(int)
			atomic.AddInt32(&sum, int32(a))
			atomic.AddInt32(&sum, int32(b))
			time.Sleep(10 * time.Millisecond)
			wg.Done()
		}, i, i)
	}
	wg.Wait()
	ericlog.Debugf("end. sum=%d", sum)
}

func TestGlobal(t *testing.T) {
	err := taskpool.Init()
	assert.Equal(t, nil, err)
	s := taskpool.GetCurrentStatus()
	assert.Equal(t, 0, s.TotalWorkerNum)
	assert.Equal(t, 0, s.IdleWorkerNum)
	assert.Equal(t, 0, s.BlockTaskNum)
	taskpool.Go(func(param ...interface{}) {
	})
	taskpool.KillIdleWorkers()
}

func TestCorner(t *testing.T) {
	_, err := taskpool.NewPool(func(option *taskpool.Option) {
		option.InitWorkerNum = -1
	})
	assert.Equal(t, taskpool.ErrTaskPool, err)

	_, err = taskpool.NewPool(func(option *taskpool.Option) {
		option.MaxWorkerNum = -1
	})
	assert.Equal(t, taskpool.ErrTaskPool, err)

	_, err = taskpool.NewPool(func(option *taskpool.Option) {
		option.InitWorkerNum = 5
		option.MaxWorkerNum = 1
	})
	assert.Equal(t, taskpool.ErrTaskPool, err)
}
