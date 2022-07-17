package env

import (
	"log"
	"os"
	"testing"

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

func TestGet(t *testing.T) {
	r := require.New(t)
	r.NotZero(os.Getenv("GOPATH"))

	r.Equal(os.Getenv("GOPATH"), Get("GOPATH", "foo"))
	r.Equal("bar", Get("IDONTEXIST", "bar"))
	r.Equal(true, Bool("ENV_DEBUG", false))
	r.Equal(8080, Int("ENV_PORT", 0))
}

func Test_MustGet(t *testing.T) {
	r := require.New(t)
	r.NotZero(os.Getenv("GOPATH"))

	v, err := MustGet("GOPATH")
	r.NoError(err)
	r.Equal(os.Getenv("GOPATH"), v)

	_, err = MustGet("IDONTEXIST")
	r.Error(err)
}

func Test_Set(t *testing.T) {
	r := require.New(t)
	r.Zero(os.Getenv("FOO"))
	_, err := MustGet("FOO")
	r.Error(err)

	Set("FOO", "foo")
	r.Equal("foo", Get("FOO", "bar"))
	// but Set should not touch the os envrionment
	r.NotEqual(os.Getenv("FOO"), "foo")
	r.Error(err)
}

func Test_MustSet(t *testing.T) {
	r := require.New(t)
	r.Zero(os.Getenv("FOO"))

	err := MustSet("FOO", "BAR")
	r.NoError(err)
	// MustSet also set underlying os environment
	r.Equal("BAR", os.Getenv("FOO"))
}

// Env files loading: automatically loaded by init()
func Test_LoadEnvLoadsEnvFile(t *testing.T) {
	r := require.New(t)
	r.Equal("root", Get("ENV_DIR", ""))
	r.Equal("none", Get("ENV_FLAVOUR", ""))
}
