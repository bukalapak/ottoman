package json_test

import (
	"strings"
	"testing"

	"github.com/bukalapak/ottoman/encoding/json"
	"github.com/go-restit/lzjson"
	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	s := `{"data":{"hello":"world","items":[1,2,3]},"meta":{"code":200}}`
	n := json.NewNode(strings.NewReader(s))

	assert.True(t, n.IsValid())
	assert.Equal(t, []byte(s), n.Bytes())

	assert.Contains(t, n.Get("data").Keys(), "hello")
	assert.Contains(t, n.Get("data").Keys(), "items")
	assert.Equal(t, []string{"code"}, n.Get("meta").Keys())
	assert.Equal(t, []string{"hello", "items"}, n.Get("data").SortedKeys())

	assert.Equal(t, 3, n.Get("data").Get("items").Len())
	assert.Equal(t, 1, n.Get("data").Get("items").GetN(0).Int())
	assert.Equal(t, "world", n.Get("data").Get("hello").String())

	assert.Equal(t, lzjson.TypeObject, n.Get("data").Type())
	assert.Equal(t, lzjson.TypeString, n.Get("data").Get("hello").Type())
	assert.Equal(t, lzjson.TypeArray, n.Get("data").Get("items").Type())
	assert.Equal(t, lzjson.TypeNumber, n.Get("meta").Get("code").Type())

	z := `{"obj":{},"arr":[],"str":"","num":0,"nil":null,"boo":true}`
	w := json.NewNode(strings.NewReader(z))

	assert.True(t, w.IsValid())
	assert.True(t, w.Get("obj").IsObject())
	assert.True(t, w.Get("obj").IsEmpty())

	assert.True(t, w.Get("arr").IsArray())
	assert.True(t, w.Get("arr").IsEmpty())

	assert.True(t, w.Get("str").IsString())
	assert.True(t, w.Get("str").IsEmpty())

	assert.True(t, w.Get("num").IsNumber())
	assert.False(t, w.Get("num").IsEmpty())

	assert.True(t, w.Get("nil").IsNull())
	assert.True(t, w.Get("boo").IsBool())

	m := map[string]int{}
	x := map[string]int{"code": 200}

	err := n.Get("meta").Unmarshal(&m)
	assert.Nil(t, err)
	assert.Equal(t, x, m)

	v := n.Get("unknown")
	assert.NotNil(t, v.Error())
}

func TestNode_invalid(t *testing.T) {
	s := "hello!"
	n := json.NewNode(strings.NewReader(s))

	assert.False(t, n.IsValid())
}
