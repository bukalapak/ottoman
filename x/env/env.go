package env

import (
	"os"
	"strconv"
)

func Get(key string) string {
	return os.Getenv(key)
}

func Set(key, value string) error {
	return os.Setenv(key, value)
}

func Unset(key string) error {
	return os.Unsetenv(key)
}

func Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

func Clear() {
	os.Clearenv()
}

func Expand(s string) string {
	return os.ExpandEnv(s)
}

func Fetch(key, def string) string {
	s, ok := Lookup(key)
	if !ok {
		return def
	}

	return s
}

func String(key string) string {
	return Get(key)
}

func Bool(key string) bool {
	s := String(key)
	if s == "yes" {
		return true
	}

	v, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}

	return v
}

func Int(key string) int {
	n, err := strconv.ParseInt(String(key), 10, 32)
	if err != nil {
		return 0
	}

	return int(n)
}

func Float64(key string) float64 {
	n, err := strconv.ParseFloat(String(key), 64)
	if err != nil {
		return 0.0
	}

	return n
}
