package header

import (
	"mime"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/subosito/httpx"
)

var (
	MediaTypeAny           = "*/*"
	MediaTypeJSON          = "application/json"
	MediaTypeMsgPack       = "application/msgpack"
	mediaTypeVendorPattern = `application/vnd.{vendor}.{version}+{format}`
	mediaTypeVendorRegex   = `application/vnd.(?P<vendor>\w+).(?P<version>v\d{1})\+(?P<format>\w+)`
)

func init() {
	mime.AddExtensionType(".json", MediaTypeJSON)
	mime.AddExtensionType(".msgpack", MediaTypeMsgPack)
}

type ContentType struct {
	MediaType  string
	Quality    float64
	Parameters map[string]string
}

type ContentTypes []ContentType

func (c ContentTypes) Len() int           { return len(c) }
func (c ContentTypes) Less(i, j int) bool { return c[i].Quality < c[j].Quality }
func (c ContentTypes) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func (c ContentTypes) ContentType() ContentType {
	if len(c) == 0 {
		return ContentType{}
	}

	if n := c.contains(MediaTypeMsgPack); n != -1 {
		return c[n]
	}

	if n := c.contains(MediaTypeJSON); n != -1 {
		return c[n]
	}

	if v, ok := c.Wildcard(); ok {
		return v
	}

	return c[0]
}

func (c ContentTypes) Wildcard() (ContentType, bool) {
	if n := c.contains(MediaTypeAny); n != -1 {
		return c[n], true
	}

	return ContentType{}, false
}

func (c ContentTypes) contains(s string) int {
	for i := range c {
		if c[i].MediaType == s {
			return i
		}
	}

	return -1
}

func (c ContentTypes) Contains(s string) bool {
	return c.contains(s) != -1
}

func (c ContentType) String() string {
	return mime.FormatMediaType(c.MediaType, c.Parameters)
}

func (c ContentType) VendorString() string {
	if s := c.vendorString(); s != mediaTypeVendorPattern {
		return s
	}

	return c.String()
}

func (c ContentType) vendorString() string {
	s := mediaTypeVendorPattern

	if c.mapContains("vendor") && c.mapContains("version") {
		if x := strings.SplitN(c.MediaType, "/", 2); len(x) == 2 {
			s = strings.Replace(s, "{format}", x[1], 1)
		}

		for k, v := range c.Parameters {
			s = strings.Replace(s, "{"+k+"}", v, 1)
		}
	}

	return s
}

func (c ContentType) mapContains(s string) bool {
	for k, v := range c.Parameters {
		if k == s && v != "" {
			return true
		}
	}

	return false
}

func NewContentType(s string, m map[string]string) ContentType {
	c := ContentType{MediaType: s, Quality: 1.0, Parameters: make(map[string]string)}

	for k, v := range m {
		c.Parameters[k] = v
	}

	for k, v := range mediaVendor(s) {
		if k == "format" {
			if z, _, err := mime.ParseMediaType(mime.TypeByExtension("." + v)); err == nil {
				c.MediaType = z
			}
		} else {
			c.Parameters[k] = v
		}
	}

	for k, v := range c.Parameters {
		if k == "q" {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				c.Quality = f
			}
		}
	}

	return c
}

func mediaVendor(s string) map[string]string {
	re := regexp.MustCompile(mediaTypeVendorRegex)
	mh := re.FindStringSubmatch(s)
	mv := make(map[string]string)

	for i, name := range re.SubexpNames() {
		if i > 0 && i <= len(mh) {
			mv[name] = mh[i]
		}
	}

	return mv
}

func ContentHeader(h http.Header, key string) ContentTypes {
	cm := ContentTypes{}
	ss := httpx.ParseList(h, key)

	for _, s := range ss {
		if m, p, err := mime.ParseMediaType(s); err == nil {
			cm = append(cm, NewContentType(m, p))
		}
	}

	sort.Sort(cm)

	return cm
}
