package proxy

import (
	"container/ring"
	"net/url"
)

type Targeter interface {
	Target() *url.URL
}

type Target struct {
	baseURL *url.URL
}

func NewTarget(baseURL *url.URL) *Target {
	return &Target{baseURL: baseURL}
}

func (t *Target) Target() *url.URL {
	return t.baseURL
}

type RingTarget struct {
	ringURL *ring.Ring
}

func NewRingTarget(baseURL []*url.URL) *RingTarget {
	return &RingTarget{ringURL: urlRing(baseURL)}
}

func (t *RingTarget) Target() *url.URL {
	u := t.ringURL.Value.(*url.URL)
	t.ringURL = t.ringURL.Next()

	return u
}

func urlRing(us []*url.URL) *ring.Ring {
	r := ring.New(len(us))

	for i := 0; i < r.Len(); i++ {
		r.Value = us[i]
		r = r.Next()
	}

	return r
}
