package env_test

import (
	"os"
	"testing"

	"github.com/bukalapak/ottoman/x/env"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	env.Set("foo", "bar")
	assert.Equal(t, "bar", os.Getenv("foo"))
	os.Clearenv()
}

func TestGet(t *testing.T) {
	os.Setenv("foo", "bar")
	assert.Equal(t, "bar", env.Get("foo"))
	os.Clearenv()
}

func TestUnset(t *testing.T) {
	os.Setenv("foo", "bar")
	env.Unset("foo")

	_, exist := os.LookupEnv("foo")

	assert.False(t, exist)
	os.Clearenv()
}

func TestLookup(t *testing.T) {
	os.Setenv("foo", "bar")

	s, ok := env.Lookup("foo")

	assert.Equal(t, "bar", s)
	assert.True(t, ok)

	os.Clearenv()
}

func TestClear(t *testing.T) {
	os.Setenv("foo", "bar")
	env.Clear()

	_, exist := os.LookupEnv("foo")

	assert.False(t, exist)
	os.Clearenv()
}

func TestExpand(t *testing.T) {
	os.Setenv("foo", "bar")
	assert.Equal(t, "Hello bar!", env.Expand("Hello ${foo}!"))
	os.Clearenv()
}

func TestFetch(t *testing.T) {
	assert.Equal(t, "boo", env.Fetch("foo", "boo"))
}

func TestFetch_valueExist(t *testing.T) {
	os.Setenv("foo", "bar")
	assert.Equal(t, "bar", env.Fetch("foo", "boo"))
	os.Clearenv()
}

func TestString(t *testing.T) {
	os.Setenv("foo", "bar")
	assert.Equal(t, "bar", env.String("foo"))
	os.Clearenv()
}

func TestBool(t *testing.T) {
	table := map[string]bool{
		"true":  true,
		"false": false,
		"yes":   true,
		"no":    false,
		"1":     true,
		"0":     false,
	}

	for k, v := range table {
		os.Setenv("foo", k)
		assert.Equal(t, v, env.Bool("foo"))
		os.Clearenv()
	}
}

func TestInt(t *testing.T) {
	table := map[string]int{
		"abc": 0,
		"5.5": 0,
		"-20": -20,
		"100": 100,
	}

	for k, v := range table {
		os.Setenv("foo", k)
		assert.Equal(t, v, env.Int("foo"))
		os.Clearenv()
	}
}

func TestFloat64(t *testing.T) {
	table := map[string]float64{
		"abcd": 0,
		"5.54": 5.54,
		"-2.2": -2.2,
		"1000": 1000,
	}

	for k, v := range table {
		os.Setenv("foo", k)
		assert.Equal(t, v, env.Float64("foo"))
		os.Clearenv()
	}
}
