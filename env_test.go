package goenv

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
	r := require.New(t)
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
	r := require.New(t)
	r.NotZero(os.Getenv("GOPATH"))
	r.True(IsSet("GOPATH"))

	r.Equal(os.Getenv("GOPATH"), Get("GOPATH", "foo"))
	r.Equal("bar", Get("IDONTEXIST", "bar"))
	r.Equal(false, Bool("ENV_DEBUG", false))
	r.Equal(false, Bool("IDONTEXIST", false))
	port, err := Int("ENV_PORT", 0)
	r.NoError(err)
	r.Equal(0, port)

	os.Setenv("Dur", "3s")
	dur, err := Duration("Dur", time.Second)
	r.NoError(err)
	r.Equal(time.Second*3, dur)
}

func TestLoad(t *testing.T) {
	r := require.New(t)
	err := Load()
	r.NoError(err)
	r.NotZero(os.Getenv("ENV_DIR"))
	r.NotZero(os.Getenv("ENV_FLAVOUR"))
	r.NotZero(os.Getenv("ENV_PORT"))
	r.NotZero(os.Getenv("ENV_DEBUG"))

	r.Equal("root", Get("ENV_DIR", ""))
	r.Equal("none", Get("ENV_FLAVOUR", ""))
	r.Equal("8080", Get("ENV_PORT", ""))
	port, err := Int("ENV_PORT", 0)
	r.NoError(err)
	r.Equal(8080, port)
	r.Equal(true, Bool("ENV_DEBUG", false))
}

func TestMarshal(t *testing.T) {
	r := require.New(t)
	m, err := Marshal()
	r.NoError(err)
	t.Log(m)
	r.NotZero(m)
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get("GOPATH", "foo")
	}
}
