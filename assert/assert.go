// package assert 提供了单元测试时的断言功能，减少一些模板代码
package assert

import (
	"bytes"
	"reflect"
)

// 单元测试中的 *testing.T 和 *testing.B 都满足该接口
type TestingT interface {
	Errorf(format string, args ...interface{})
}

type tHelper interface {
	Helper()
}

func Equal(t TestingT, expected interface{}, actual interface{}, msg ...string) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !equal(expected, actual) {
		t.Errorf("%s expected=%+v, actual=%+v", msg, expected, actual)
	}
	return
}

// 比如有时我们需要对 error 类型不等于 nil 做断言，但是我们并不关心 error 的具体值是什么
func IsNotNil(t TestingT, actual interface{}, msg ...string) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if isNil(actual) {
		t.Errorf("%s expected not nil, but actual=%+v", msg, actual)
	}
	return
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
