package ericlog_test

import (
	"encoding/hex"
	"fmt"
	originLog "log"
	"os"
	"sync"
	"testing"
	"time"

	ericlog "github.com/ericxiao417/ericool/ericlog"

	"github.com/ericxiao417/ericool/fake"

	"github.com/ericxiao417/ericool/assert"
)

func TestLogger(t *testing.T) {
	l, err := ericlog.New(func(option *ericlog.Option) {
		option.Level = ericlog.LevelInfo
		option.Filename = "/tmp/ericlogtest/aaa.log"
		option.IsToStdout = true
		option.IsRotateDaily = true
	})
	assert.Equal(t, nil, err)
	buf := []byte("1234567890987654321")
	l.Error(hex.Dump(buf))
	l.Debugf("l test msg by Debug%s", "f")
	l.Infof("l test msg by Info%s", "f")
	l.Warnf("l test msg by Warn%s", "f")
	l.Errorf("l test msg by Error%s", "f")
	l.Debug("l test msg by Debug")
	l.Info("l test msg by Info")
	l.Warn("l test msg by Warn")
	l.Error("l test msg by Error")
	l.Output(2, "l test msg by Output")
	l.Out(ericlog.LevelInfo, 1, "l test msg by Out")
	l.Print("l test msg by Print")
	l.Printf("l test msg by Print%s", "f")
	l.Println("l test msg by Print")
}

func TestGlobal(t *testing.T) {
	buf := []byte("1234567890987654321")
	ericlog.Error(hex.Dump(buf))
	ericlog.Debugf("g test msg by Debug%s", "f")
	ericlog.Infof("g test msg by Info%s", "f")
	ericlog.Warnf("g test msg by Warn%s", "f")
	ericlog.Errorf("g test msg by Error%s", "f")
	ericlog.Debug("g test msg by Debug")
	ericlog.Info("g test msg by Info")
	ericlog.Warn("g test msg by Warn")
	ericlog.Error("g test msg by Error")

	err := ericlog.Init(func(option *ericlog.Option) {
		option.Level = ericlog.LevelInfo
		option.Filename = "/tmp/ericlogtest/bbb.log"
		option.IsToStdout = true

	})
	assert.Equal(t, nil, err)
	ericlog.Debugf("gc test msg by Debug%s", "f")
	ericlog.Infof("gc test msg by Info%s", "f")
	ericlog.Warnf("gc test msg by Warn%s", "f")
	ericlog.Errorf("gc test msg by Error%s", "f")
	ericlog.Debug("gc test msg by Debug")
	ericlog.Info("gc test msg by Info")
	ericlog.Warn("gc test msg by Warn")
	ericlog.Error("gc test msg by Error")
	ericlog.Output(2, "gc test msg by Output")
	ericlog.Out(ericlog.LevelInfo, 2, "gc test msg by Out")
	ericlog.Print("gc test msg by Print")
	ericlog.Printf("gc test msg by Print%s", "f")
	ericlog.Println("gc test msg by Print")
	ericlog.Sync()
}

func TestNew(t *testing.T) {
	var (
		l   ericlog.Logger
		err error
	)
	l, err = ericlog.New(func(option *ericlog.Option) {
		option.Level = ericlog.LevelPanic + 1
	})
	assert.Equal(t, nil, l)
	assert.Equal(t, ericlog.ErrLog, err)

	l, err = ericlog.New(func(option *ericlog.Option) {
		option.AssertBehavior = ericlog.AssertPanic + 1
	})
	assert.Equal(t, nil, l)
	assert.Equal(t, ericlog.ErrLog, err)

	l, err = ericlog.New(func(option *ericlog.Option) {
		option.Filename = "/tmp"
	})
	assert.Equal(t, nil, l)
	assert.IsNotNil(t, err)

	l, err = ericlog.New(func(option *ericlog.Option) {
		option.Filename = "./log_test.go/111"
	})
	assert.Equal(t, nil, l)
	assert.IsNotNil(t, err)
}

func TestRotate(t *testing.T) {
	err := ericlog.Init(func(option *ericlog.Option) {
		option.Level = ericlog.LevelInfo
		option.Filename = "/tmp/ericlogtest/ccc.log"
		option.IsToStdout = false
		option.IsRotateDaily = true

	})
	assert.Equal(t, nil, err)
	ericlog.Info("aaa")
	fake.WithFakeTimeNow(func() time.Time {
		return time.Now().Add(48 * time.Hour)
	}, func() {
		ericlog.Info("bbb")
	})
}

func TestPanic(t *testing.T) {
	fake.WithRecover(func() {
		ericlog.Panic("aaa")
	})
	fake.WithRecover(func() {
		ericlog.Panicf("%s", "bbb")
	})
	fake.WithRecover(func() {
		ericlog.Panicln("aaa")
	})
	fake.WithRecover(func() {
		l, err := ericlog.New()
		assert.Equal(t, nil, err)
		l.Panic("aaa")
	})
	fake.WithRecover(func() {
		l, err := ericlog.New()
		assert.Equal(t, nil, err)
		l.Panicf("%s", "bbb")
	})
	fake.WithRecover(func() {
		l, err := ericlog.New()
		assert.Equal(t, nil, err)
		l.Panicln("aaa")
	})
}

func TestFatal(t *testing.T) {
	var er fake.ExitResult

	er = fake.WithFakeOSExit(func() {
		ericlog.Fatal("Fatal")
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)

	er = fake.WithFakeOSExit(func() {
		ericlog.Fatalf("Fatalf%s", ".")
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)

	er = fake.WithFakeOSExit(func() {
		ericlog.Fatalln("Fatalln")
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)

	logger, err := ericlog.New(func(option *ericlog.Option) {
		option.Level = ericlog.LevelInfo
	})
	assert.IsNotNil(t, logger)
	assert.Equal(t, nil, err)
	er = fake.WithFakeOSExit(func() {
		logger.Fatal("Fatal")
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)

	er = fake.WithFakeOSExit(func() {
		logger.Fatalf("Fatalf%s", ".")
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)

	er = fake.WithFakeOSExit(func() {
		logger.Fatalln("Fatalln")
	})
	assert.Equal(t, true, er.HasExit)
	assert.Equal(t, 1, er.ExitCode)
}

func TestAssert(t *testing.T) {
	// 成功
	ericlog.Assert(nil, nil)
	ericlog.Assert(nil, nil)
	ericlog.Assert(nil, nil)
	ericlog.Assert(1, 1)
	ericlog.Assert("aaa", "aaa")
	var ch chan struct{}
	ericlog.Assert(nil, ch)
	var m map[string]string
	ericlog.Assert(nil, m)
	var p *int
	ericlog.Assert(nil, p)
	var i interface{}
	ericlog.Assert(nil, i)
	var b []byte
	ericlog.Assert(nil, b)

	ericlog.Assert([]byte{}, []byte{})
	ericlog.Assert([]byte{0, 1, 2}, []byte{0, 1, 2})

	// 失败
	_ = ericlog.Init(func(option *ericlog.Option) {
		option.AssertBehavior = ericlog.AssertError
	})
	ericlog.Assert(nil, 1)

	_ = ericlog.Init(func(option *ericlog.Option) {
		option.AssertBehavior = ericlog.AssertFatal
	})
	err := fake.WithFakeOSExit(func() {
		ericlog.Assert(nil, 1)
	})
	assert.Equal(t, true, err.HasExit)
	assert.Equal(t, 1, err.ExitCode)

	_ = ericlog.Init(func(option *ericlog.Option) {
		option.AssertBehavior = ericlog.AssertPanic
	})
	fake.WithRecover(func() {
		ericlog.Assert([]byte{}, "aaa")
	})

	l, _ := ericlog.New()
	l.Assert(nil, 1)

	l, _ = ericlog.New(func(option *ericlog.Option) {
		option.AssertBehavior = ericlog.AssertFatal
	})
	err = fake.WithFakeOSExit(func() {
		l.Assert(nil, 1)
	})
	assert.Equal(t, true, err.HasExit)
	assert.Equal(t, 1, err.ExitCode)

	l, _ = ericlog.New(func(option *ericlog.Option) {
		option.AssertBehavior = ericlog.AssertPanic
	})
	fake.WithRecover(func() {
		l.Assert([]byte{}, "aaa")
	})
}

func TestLogger_WithPrefix(t *testing.T) {
	im := 4
	jm := 4
	var wg sync.WaitGroup
	wg.Add(im * jm)
	ericlog.Debug(">")
	for i := 0; i != im; i++ {
		go func(ii int) {
			for j := 0; j != jm; j++ {
				s := fmt.Sprintf("%d", ii)
				l := ericlog.WithPrefix("log_test")
				l.Info(j)
				ll := l.WithPrefix("TestLogger_WithPrefix")
				ll.Info(j)
				lll := ll.WithPrefix(s)
				lll.Info(j)
				wg.Done()
			}
		}(i)
	}
	ericlog.Debug("<")
	wg.Wait()
}

func Benchmarkericlog(b *testing.B) {
	b.ReportAllocs()

	err := ericlog.Init(func(option *ericlog.Option) {
		option.Level = ericlog.LevelInfo
		option.Filename = "/dev/null"
		option.IsToStdout = false
		option.IsRotateDaily = false
	})
	assert.Equal(b, nil, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ericlog.Infof("hello %s %d", "world", i)
		ericlog.Info("Info")
	}
}

func BenchmarkOriginLog(b *testing.B) {
	b.ReportAllocs()

	fp, err := os.Create("/dev/null")
	assert.Equal(b, nil, err)
	originLog.SetOutput(fp)
	originLog.SetFlags(originLog.Ldate | originLog.Ltime | originLog.Lmicroseconds | originLog.Lshortfile)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		originLog.Printf("hello %s %d\n", "world", i)
		originLog.Println("Info")
	}
}
