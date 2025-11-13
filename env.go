// Package goenv provides functions to manage environment variables
// from .env files and retrieve typed values from the environment.
package goenv

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

// IsSet returns if the given env key is set.
// remember ENV must be a non-empty. All empty
// values are considered unset.
func IsSet(key string) bool {
	return Get(key, "") != ""
}

// Get a value from the ENV. If it doesn't exist the
// default value will be returned.
func Get(key string, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return fastTrim(v)
	}
	return defaultValue
}

// String returns the string value represented by the string.
func String(key string, defaultValue string) string {
	return Get(key, defaultValue)
}

// Bool returns the boolean value represented by the string.
func Bool(key string, defaultValue bool) bool {
	val := Get(key, "")
	if val == "true" ||
		val == "1" ||
		val == "t" ||
		val == "T" ||
		val == "TRUE" ||
		val == "True" {
		return true
	}
	return defaultValue
}

// Int returns the integer value represented by the string.
// It supports all integer types (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64).
func Int[T constraints.Integer](key string, defaultValue T) (T, error) {
	v := Get(key, "")
	if v == "" {
		return defaultValue, nil
	}

	var result T
	var err error

	switch any(result).(type) {
	case int:
		val, e := strconv.Atoi(v)
		result = any(val).(T)
		err = e
	case int8:
		val, e := strconv.ParseInt(v, 10, 8)
		result = any(int8(val)).(T)
		err = e
	case int16:
		val, e := strconv.ParseInt(v, 10, 16)
		result = any(int16(val)).(T)
		err = e
	case int32:
		val, e := strconv.ParseInt(v, 10, 32)
		result = any(int32(val)).(T)
		err = e
	case int64:
		val, e := strconv.ParseInt(v, 10, 64)
		result = any(val).(T)
		err = e
	case uint:
		val, e := strconv.ParseUint(v, 10, 0)
		result = any(uint(val)).(T)
		err = e
	case uint8:
		val, e := strconv.ParseUint(v, 10, 8)
		result = any(uint8(val)).(T)
		err = e
	case uint16:
		val, e := strconv.ParseUint(v, 10, 16)
		result = any(uint16(val)).(T)
		err = e
	case uint32:
		val, e := strconv.ParseUint(v, 10, 32)
		result = any(uint32(val)).(T)
		err = e
	case uint64:
		val, e := strconv.ParseUint(v, 10, 64)
		result = any(val).(T)
		err = e
	default:
		err = fmt.Errorf("unsupported integer type")
	}

	return result, err
}

// Duration returns a parsed time.Duration if found in
// the environment value, returns the default value duration
// otherwise.
func Duration(key string, defaultValue time.Duration) (time.Duration, error) {
	v := Get(key, "")
	if v == "" {
		return defaultValue, nil
	}
	return time.ParseDuration(v)
}

// Load loads environment variables from .env file(s).
// It searches for .env file in the current directory by default.
// If multiple files are provided, they will be loaded in order.
// Returns an error if any file cannot be loaded.
func Load(filenames ...string) (err error) {
	filenames = filenamesOrDefault(filenames)

	for _, filename := range filenames {
		err = loadFile(filename, false)
		if err != nil {
			return // return early on a spazout
		}
	}
	return
}

// Marshal outputs the given environment as a dotenv-formatted environment file.
// Each line is in the format: KEY="VALUE" where VALUE is backslash-escaped.
func Marshal() (string, error) {
	var envMap = map[string]string{}

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		envMap[pair[0]] = fastTrim(os.Getenv(pair[0]))
	}
	lines := make([]string, 0, len(envMap))
	for k, v := range envMap {
		if d, err := strconv.Atoi(v); err == nil {
			lines = append(lines, fmt.Sprintf(`%s=%d`, k, d))
		} else {
			lines = append(lines, fmt.Sprintf(`%s="%s"`, k, doubleQuoteEscape(v)))
		}
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n"), nil
}

const doubleQuoteSpecialChars = "\\\n\r\"!$`"

func doubleQuoteEscape(line string) string {
	for _, c := range doubleQuoteSpecialChars {
		toReplace := "\\" + string(c)
		if c == '\n' {
			toReplace = `\n`
		}
		if c == '\r' {
			toReplace = `\r`
		}
		line = strings.Replace(line, string(c), toReplace, -1)
	}
	return line
}
func filenamesOrDefault(filenames []string) []string {
	if len(filenames) == 0 {
		return []string{".env"}
	}
	return filenames
}

func loadFile(filename string, overload bool) error {
	envMap, err := readFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] || overload {
			_ = os.Setenv(key, value)
		}
	}

	return nil
}

func readFile(filename string) (envMap map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}
	envMap = map[string]string{}
	err = parseBytes(buf.Bytes(), envMap)
	return
}

func fastTrim(s string) string {
	if s == "" {
		return s
	}

	start := 0
	end := len(s)
	for start < end {
		if s[start] != ' ' {
			break
		}
		start++
	}
	for end > start {
		if s[end-1] != ' ' {
			break
		}
		end--
	}
	if start == 0 && end == len(s) {
		return s
	}
	return s[start:end]
}
