package goenv

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var mu = &sync.RWMutex{}
var env = map[string]string{}

func init() {
	Load()
	loadSystemEnv()
}

// Load the ENV variables to the env map
func loadSystemEnv() {
	mu.Lock()
	defer mu.Unlock()

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		env[pair[0]] = os.Getenv(pair[0])
	}
}

// Reload the ENV variables. Useful if
// an external ENV manager has been used
func Reload() {
	env = map[string]string{}
	loadSystemEnv()
}

// Load .env files. Files will be loaded in the same order that are received.
// Redefined vars will override previously existing values.
// IE: envy.Load(".env", "test_env/.env") will result in DIR=test_env
// If no arg passed, it will try to load a .env file.
func Load(files ...string) error {

	// If no files received, load the default one
	if len(files) == 0 {
		err := godotenv.Load()
		if err == nil {
			Reload()
		}
		return err
	}

	// We received a list of files
	for _, file := range files {

		// Check if it exists or we can access
		if _, err := os.Stat(file); err != nil {
			// It does not exist or we can not access.
			// Return and stop loading
			return err
		}

		// It exists and we have permission. Load it
		if err := godotenv.Load(file); err != nil {
			return err
		}

		// Reload the env so all new changes are noticed
		Reload()

	}
	return nil
}

// Get a value from the ENV. If it doesn't exist the
// default value will be returned.
func Get(key string, value string) string {
	mu.RLock()
	defer mu.RUnlock()
	if v, ok := env[key]; ok {
		return v
	}
	return value
}

// MustSet the value into the underlying ENV, as well as envy.
// This may return an error if there is a problem setting the
// underlying ENV value.
func MustSet(key string, value string) error {
	mu.Lock()
	defer mu.Unlock()
	err := os.Setenv(key, value)
	if err != nil {
		return err
	}
	env[key] = value
	return nil
}

// MustGet get a value from the ENV. If it doesn't exist
// an error will be returned
func MustGet(key string) (string, error) {
	mu.RLock()
	defer mu.RUnlock()
	if v, ok := env[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("could not find ENV var with %s", key)
}

// Set a value into the ENV. This is NOT permanent. It will
// only affect values accessed through envy.
func Set(key string, value string) {
	mu.Lock()
	defer mu.Unlock()
	env[key] = value
}

//Bool returns the boolean value represented by the string.
func Bool(key string, value bool) bool {
	mu.RLock()
	defer mu.RUnlock()
	if v, ok := env[key]; ok {
		return strings.ToLower(v) == "true"
	}
	return value
}

//Int returns the integer value represented by the string.
func Int(key string, value int) int {
	mu.RLock()
	defer mu.RUnlock()
	if v, ok := env[key]; ok {
		i, _ := strconv.Atoi(v)
		return i
	}
	return value
}
