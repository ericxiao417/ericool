package taskpool_test

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/ericxiao417/ericool/taskpool"
)

// 并发计算0+1+2+...+1000
// 演示怎么向协程池中添加带参数的函数任务
func ExampleNewPool() {
	pool, _ := taskpool.NewPool(func(option *taskpool.Option) {
		// 限制最大并发数
		option.MaxWorkerNum = 16
	})
	var sum int32
	var wg sync.WaitGroup
	n := 1000
	wg.Add(n)
	for i := 0; i < n; i++ {
		pool.Go(func(param ...interface{}) {
			ii := param[0].(int)
			atomic.AddInt32(&sum, int32(ii))
			wg.Done()
		}, i)
	}
	wg.Wait()
	fmt.Println(sum)
	// Output:
	// 499500
}
