package proxy_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/bukalapak/ottoman/proxy"
	"github.com/stretchr/testify/assert"
)

func TestTarget(t *testing.T) {
	u, _ := url.Parse("http://api.example.com/")
	x := proxy.NewTarget(u)
	z := x.Target()

	assert.Equal(t, "http", z.Scheme)
	assert.Equal(t, "api.example.com", z.Host)
}

func TestRingTarget(t *testing.T) {
	us := make([]*url.URL, 3)
	ss := []string{
		"http://s1.example.com/",
		"http://s2.example.com/",
		"http://s3.example.com/",
	}

	for i, s := range ss {
		u, _ := url.Parse(s)
		us[i] = u
	}

	x := proxy.NewRingTarget(us)

	for i := 1; i <= 3; i++ {
		z := x.Target()
		assert.Equal(t, "http", z.Scheme)
		assert.Equal(t, fmt.Sprintf("s%d.example.com", i), z.Host)
	}
}
