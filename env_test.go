package goenv

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dnsoa/go/assert"
)

var _ = func() error {
	str := `
# This is a comment
# We can use equal or colon notation
ENV_DIR: root
ENV_FLAVOUR: none
ENV_PORT: 8080
ENV_DEBUG: true
`
	f, err := os.Create(".env")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(str)
	return nil
}()

func TestTrim(t *testing.T) {
	r := assert.New(t)
	r.Equal("", fastTrim(""))
	r.Equal("foo", fastTrim("foo"))
	r.Equal("foo", fastTrim(" foo "))
	r.Equal("foo", fastTrim("foo "))
	r.Equal("foo", fastTrim(" foo"))
}

func BenchmarkTrim(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.TrimSpace(" foo ")
	}
}

func TestGet(t *testing.T) {
	r := assert.New(t)
	r.NotEmpty(os.Getenv("GOPATH"))
	r.True(IsSet("GOPATH"))

	r.Equal(os.Getenv("GOPATH"), Get("GOPATH", "foo"))
	r.Equal("bar", Get("IDONTEXIST", "bar"))
	r.Equal(false, Bool("ENV_DEBUG", false))
	r.Equal(false, Bool("IDONTEXIST", false))
	port, err := Int[int]("ENV_PORT", 0)
	r.NoError(err)
	r.Equal(0, port)

	os.Setenv("Dur", "3s")
	dur, err := Duration("Dur", time.Second)
	r.NoError(err)
	r.Equal(time.Second*3, dur)
}

func TestLoad(t *testing.T) {
	r := assert.New(t)
	err := Load()
	r.NoError(err)
	r.NotEmpty(os.Getenv("ENV_DIR"))
	r.NotEmpty(os.Getenv("ENV_FLAVOUR"))
	r.NotEmpty(os.Getenv("ENV_PORT"))
	r.NotEmpty(os.Getenv("ENV_DEBUG"))

	r.Equal("root", Get("ENV_DIR", ""))
	r.Equal("none", Get("ENV_FLAVOUR", ""))
	r.Equal("8080", Get("ENV_PORT", ""))
	port, err := Int[int]("ENV_PORT", 0)
	r.NoError(err)
	r.Equal(8080, port)
	r.Equal(true, Bool("ENV_DEBUG", false))
}

func TestMarshal(t *testing.T) {
	r := assert.New(t)
	m, err := Marshal()
	r.NoError(err)
	t.Log(m)
	r.NotEmpty(m)
}

func TestGenericInt(t *testing.T) {
	r := assert.New(t)

	// 设置测试环境变量
	os.Setenv("TEST_INT", "42")
	os.Setenv("TEST_INT8", "127")
	os.Setenv("TEST_INT16", "32767")
	os.Setenv("TEST_INT32", "2147483647")
	os.Setenv("TEST_INT64", "9223372036854775807")
	os.Setenv("TEST_UINT", "42")
	os.Setenv("TEST_UINT8", "255")
	os.Setenv("TEST_UINT16", "65535")
	os.Setenv("TEST_UINT32", "4294967295")
	os.Setenv("TEST_UINT64", "18446744073709551615")

	// 测试 int
	val, err := Int[int]("TEST_INT", 0)
	r.NoError(err)
	r.Equal(42, val)

	// 测试 int8
	val8, err := Int[int8]("TEST_INT8", 0)
	r.NoError(err)
	r.Equal(int8(127), val8)

	// 测试 int16
	val16, err := Int[int16]("TEST_INT16", 0)
	r.NoError(err)
	r.Equal(int16(32767), val16)

	// 测试 int32
	val32, err := Int[int32]("TEST_INT32", 0)
	r.NoError(err)
	r.Equal(int32(2147483647), val32)

	// 测试 int64
	val64, err := Int[int64]("TEST_INT64", 0)
	r.NoError(err)
	r.Equal(int64(9223372036854775807), val64)

	// 测试 uint
	uval, err := Int[uint]("TEST_UINT", 0)
	r.NoError(err)
	r.Equal(uint(42), uval)

	// 测试 uint8
	uval8, err := Int[uint8]("TEST_UINT8", 0)
	r.NoError(err)
	r.Equal(uint8(255), uval8)

	// 测试 uint16
	uval16, err := Int[uint16]("TEST_UINT16", 0)
	r.NoError(err)
	r.Equal(uint16(65535), uval16)

	// 测试 uint32
	uval32, err := Int[uint32]("TEST_UINT32", 0)
	r.NoError(err)
	r.Equal(uint32(4294967295), uval32)

	// 测试 uint64
	uval64, err := Int[uint64]("TEST_UINT64", 0)
	r.NoError(err)
	r.Equal(uint64(18446744073709551615), uval64)

	// 测试默认值
	defVal, err := Int[int]("NONEXISTENT", 999)
	r.NoError(err)
	r.Equal(999, defVal)

	// 测试错误情况 - 设置一个无效值然后尝试解析
	os.Setenv("TEST_INVALID", "invalid")
	_, err = Int[int]("TEST_INVALID", 0)
	r.Error(err)
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get("GOPATH", "foo")
	}
}
