package cache_test

import (
	"testing"

	"github.com/bukalapak/ottoman/cache"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	data := []string{
		"foo/bar",
		"api:foo/bar",
		"bar:foo/bar",
	}

	for i := range data {
		key := cache.Normalize(data[i], "")
		assert.Equal(t, "foo/bar", key)

		key = cache.Normalize(data[i], "api")
		assert.Equal(t, "api:foo/bar", key)

		key = cache.Normalize(data[i], "foo")
		assert.Equal(t, "foo:foo/bar", key)
	}
}

func TestNormalizeMulti(t *testing.T) {
	data := []string{
		"foo/bar",
		"api:foo/bar",
		"bar:foo/bar",
	}

	results := cache.NormalizeMulti(data, "")
	assert.Equal(t, []string{"foo/bar", "foo/bar", "foo/bar"}, results)

	results = cache.NormalizeMulti(data, "api")
	assert.Equal(t, []string{"api:foo/bar", "api:foo/bar", "api:foo/bar"}, results)

	results = cache.NormalizeMulti(data, "foo")
	assert.Equal(t, []string{"foo:foo/bar", "foo:foo/bar", "foo:foo/bar"}, results)
}
