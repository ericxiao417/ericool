package ericlog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/ericxiao417/ericool/fake"
)

var _ Logger = new(logger)

const (
	levelDebugString = "DEBUG "
	levelInfoString  = " INFO "
	levelWarnString  = " WARN "
	levelErrorString = "ERROR "
	levelFatalString = "FATAL "
	levelPanicString = "PANIC "

	levelDebugColorString = "\033[22;37mDEBUG\033[0m "
	levelInfoColorString  = "\033[22;36m INFO\033[0m "
	levelWarnColorString  = "\033[22;33m WARN\033[0m "
	levelErrorColorString = "\033[22;31mERROR\033[0m "
	levelFatalColorString = "\033[22;31mFATAL\033[0m " // 颜色和 error 级别一样
	levelPanicColorString = "\033[22;31mPANIC\033[0m " // 颜色和 error 级别一样
)

var (
	levelToString = map[Level]string{
		LevelDebug: levelDebugString,
		LevelInfo:  levelInfoString,
		LevelWarn:  levelWarnString,
		LevelError: levelErrorString,
		LevelFatal: levelFatalString,
		LevelPanic: levelPanicString,
	}
	levelToColorString = map[Level]string{
		LevelDebug: levelDebugColorString,
		LevelInfo:  levelInfoColorString,
		LevelWarn:  levelWarnColorString,
		LevelError: levelErrorColorString,
		LevelFatal: levelFatalColorString,
		LevelPanic: levelPanicColorString,
	}
)

type logger struct {
	prefixs []string
	core    *core
}

type core struct {
	option Option

	m             sync.Mutex
	fp            *os.File
	console       *os.File
	buf           bytes.Buffer
	currRoundTime time.Time
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.Out(LevelDebug, 2, fmt.Sprintf(format, v...))
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.Out(LevelInfo, 2, fmt.Sprintf(format, v...))
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.Out(LevelWarn, 2, fmt.Sprintf(format, v...))
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.Out(LevelError, 2, fmt.Sprintf(format, v...))
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Out(LevelFatal, 2, fmt.Sprintf(format, v...))
	fake.OS_Exit(1)
}

func (l *logger) Panicf(format string, v ...interface{}) {
	l.Out(LevelPanic, 2, fmt.Sprintf(format, v...))
	panic(fmt.Sprintf(format, v...))
}

func (l *logger) Debug(v ...interface{}) {
	l.Out(LevelDebug, 2, fmt.Sprint(v...))
}

func (l *logger) Info(v ...interface{}) {
	l.Out(LevelInfo, 2, fmt.Sprint(v...))
}

func (l *logger) Warn(v ...interface{}) {
	l.Out(LevelWarn, 2, fmt.Sprint(v...))
}

func (l *logger) Error(v ...interface{}) {
	l.Out(LevelError, 2, fmt.Sprint(v...))
}

func (l *logger) Fatal(v ...interface{}) {
	l.Out(LevelFatal, 2, fmt.Sprint(v...))
	fake.OS_Exit(1)
}

func (l *logger) Panic(v ...interface{}) {
	l.Out(LevelPanic, 2, fmt.Sprint(v...))
	panic(fmt.Sprint(v...))
}

func (l *logger) Output(calldepth int, s string) error {
	l.Out(LevelInfo, calldepth, s)
	return nil
}

func (l *logger) Print(v ...interface{}) {
	l.Out(LevelInfo, 2, fmt.Sprint(v...))
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.Out(LevelInfo, 2, fmt.Sprintf(format, v...))
}

func (l *logger) Println(v ...interface{}) {
	l.Out(LevelInfo, 2, fmt.Sprint(v...))
}

func (l *logger) Fatalln(v ...interface{}) {
	l.Out(LevelInfo, 2, fmt.Sprint(v...))
	fake.OS_Exit(1)
}

func (l *logger) Panicln(v ...interface{}) {
	l.Out(LevelInfo, 2, fmt.Sprint(v...))
	panic(fmt.Sprint(v...))
}

func (l *logger) Assert(expected interface{}, actual interface{}) {
	if !equal(expected, actual) {
		err := fmt.Sprintf("assert failed. excepted=%+v, but actual=%+v", expected, actual)
		switch l.core.option.AssertBehavior {
		case AssertError:
			l.Out(LevelError, 2, err)
		case AssertFatal:
			l.Out(LevelFatal, 2, err)
			fake.OS_Exit(1)
		case AssertPanic:
			l.Out(LevelPanic, 2, err)
			panic(err)
		}
	}
}

func (l *logger) Out(level Level, calldepth int, s string) {
	if l.core.option.Level > level {
		return
	}

	now := fake.Time_Now()

	var file string
	var line int
	if l.core.option.ShortFileFlag {
		_, file, line, _ = runtime.Caller(calldepth)

	}

	l.core.m.Lock()

	l.core.buf.Reset()

	writeTime(&l.core.buf, now)

	if l.core.console != nil {
		l.core.buf.WriteString(levelToColorString[level])
	} else {
		l.core.buf.WriteString(levelToString[level])
	}

	if l.prefixs != nil {
		for _, s := range l.prefixs {
			l.core.buf.WriteString("[")
			l.core.buf.WriteString(s)
			l.core.buf.WriteString("] ")
		}
	}

	l.core.buf.WriteString(s)

	if file != "" && line > 0 {
		l.core.buf.WriteString(" - ")

		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short

		l.core.buf.WriteString(file)
		l.core.buf.WriteByte(':')
		itoa(&l.core.buf, line, -1)
	}

	if l.core.buf.Len() == 0 || l.core.buf.Bytes()[l.core.buf.Len()-1] != '\n' {
		l.core.buf.WriteByte('\n')
	}

	// 输出至控制台
	if l.core.console != nil {
		_, _ = l.core.console.Write(l.core.buf.Bytes())
		if level == LevelFatal || level == LevelPanic {
			_ = l.core.console.Sync()
		}
	}

	// 输出至日志文件
	if l.core.fp != nil {
		if l.core.option.IsRotateDaily && now.Day() != l.core.currRoundTime.Day() {
			backupName := l.core.option.Filename + "." + l.core.currRoundTime.Format("20060102")
			if err := os.Rename(l.core.option.Filename, backupName); err == nil {
				_ = l.core.fp.Close()
				l.core.fp, _ = os.Create(l.core.option.Filename)
			}
			l.core.currRoundTime = now
		}
		_, _ = l.core.fp.Write(l.core.buf.Bytes())
		if level == LevelFatal || level == LevelPanic {
			_ = l.core.fp.Sync()
		}
	}

	l.core.m.Unlock()
}

func (l *logger) Sync() {
	l.core.m.Lock()
	defer l.core.m.Unlock()

	if l.core.console != nil {
		_ = l.core.console.Sync()
	}
	if l.core.fp != nil {
		_ = l.core.fp.Sync()
	}
}

func (l *logger) WithPrefix(s string) Logger {
	var prefixs []string
	if l.prefixs != nil {
		prefixs = make([]string, len(l.prefixs))
		copy(prefixs, l.prefixs)
	}
	prefixs = append(prefixs, s)
	ll := &logger{
		prefixs: prefixs,
		core:    l.core,
	}
	return ll
}

func newLogger(modOptions ...ModOption) (*logger, error) {
	var err error

	l := &logger{
		core: &core{
			currRoundTime: time.Now(),
		},
	}
	l.core.option = defaultOption

	for _, fn := range modOptions {
		fn(&l.core.option)
	}

	if err := validate(l.core.option); err != nil {
		return nil, err
	}
	if l.core.option.Filename != "" {
		dir := filepath.Dir(l.core.option.Filename)
		if err = os.MkdirAll(dir, 0777); err != nil {
			return nil, err
		}
		if l.core.fp, err = os.OpenFile(l.core.option.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
			return nil, err
		}
	}
	if l.core.option.IsToStdout {
		l.core.console = os.Stdout
	}

	return l, nil
}

func validate(option Option) error {
	if option.Level < LevelDebug || option.Level > LevelPanic {
		return ErrLog
	}
	if option.AssertBehavior < AssertError || option.AssertBehavior > AssertPanic {
		return ErrLog
	}
	return nil
}
func isNil(actual interface{}) bool {
	if actual == nil {
		return true
	}
	v := reflect.ValueOf(actual)
	k := v.Kind()
	if k == reflect.Chan || k == reflect.Map || k == reflect.Ptr || k == reflect.Interface || k == reflect.Slice {
		return v.IsNil()
	}
	return false
}

func equal(expected, actual interface{}) bool {
	if expected == nil {
		return isNil(actual)
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	return bytes.Equal(exp, act)
}

func writeTime(buf *bytes.Buffer, t time.Time) {
	year, month, day := t.Date()
	itoa(buf, year, 4)
	buf.WriteByte('/')
	itoa(buf, int(month), 2)
	buf.WriteByte('/')
	itoa(buf, day, 2)
	buf.WriteByte(' ')

	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	buf.WriteByte(':')
	itoa(buf, min, 2)
	buf.WriteByte(':')
	itoa(buf, sec, 2)
	buf.WriteByte('.')
	itoa(buf, t.Nanosecond()/1e3, 6)
	buf.WriteByte(' ')
}

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// @NOTICE 该函数拷贝自 Go 标准库 /src/log/log.go: func itoa(buf *[]byte, i int, wid int)
// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *bytes.Buffer, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	buf.Write(b[bp:])
}
