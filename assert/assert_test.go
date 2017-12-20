package assert_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ericxiao417/ericool/assert"
)

// 大部分时候 TestingT interface 的实例为单元测试中的 *testing.T 和 *testing.B
// MockTestingT 是为了对自身做单元测试
type MockTestingT struct {
}

func (mtt MockTestingT) Errorf(format string, args ...interface{}) {
	_ = fmt.Errorf(format, args...)
}
func TestEqual(t *testing.T) {
	// 测试Equal
	assert.Equal(t, nil, nil)
	assert.Equal(t, nil, nil, "fxxk.")
	assert.Equal(t, 1, 1)
	assert.Equal(t, "aaa", "aaa")
	var ch chan struct{}
	assert.Equal(t, nil, ch)
	var m map[string]string
	assert.Equal(t, nil, m)
	var p *int
	assert.Equal(t, nil, p)
	var i interface{}
	assert.Equal(t, nil, i)
	var b []byte
	assert.Equal(t, nil, b)

	assert.Equal(t, []byte{}, []byte{})
	assert.Equal(t, []byte{0, 1, 2}, []byte{0, 1, 2})

	// 测试Equal失败
	var mtt MockTestingT
	assert.Equal(mtt, nil, 1)
	assert.Equal(mtt, []byte{}, "aaa")
	assert.Equal(mtt, nil, errors.New("mock error"))
}

func TestIsNotNil(t *testing.T) {
	assert.IsNotNil(t, 1)

	// 测试IsNotNil失败
	var mtt MockTestingT
	assert.IsNotNil(mtt, nil)
}
