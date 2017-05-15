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

	m := map[string]int{}
	x := map[string]int{"code": 200}

	err := n.Get("meta").Unmarshal(&m)
	assert.Nil(t, err)
	assert.Equal(t, x, m)

	v := n.Get("unknown")
	assert.NotNil(t, v.Error())
}
