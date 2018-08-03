package json

import (
	"io"
	"sort"

	"github.com/go-restit/lzjson"
)

// Node is an interface for all JSON nodes.
type Node struct {
	node lzjson.Node
}

// Bytes returns the raw JSON string in []byte.
func (n *Node) Bytes() []byte {
	return n.node.Raw()
}

// Error returns the JSON parse error, if any.
func (n *Node) Error() error {
	return n.node.ParseError()
}

// Int unmarshal the JSON into int.
func (n *Node) Int() int {
	return n.node.Int()
}

// Keys gets an array object's keys, or nil if not an object.
func (n *Node) Keys() []string {
	return n.node.GetKeys()
}

// SortedKeys gets an array object's keys in alphabetical order
func (n *Node) SortedKeys() []string {
	ss := n.Keys()
	sort.Strings(ss)

	return ss
}

// Len gets the length of the value.
// Only works with Array and String value type.
func (n *Node) Len() int {
	return n.node.Len()
}

// String unmarshal the JSON into string.
func (n *Node) String() string {
	return n.node.String()
}

// Type returns the lzjson.Type of the containing JSON value.
func (n *Node) Type() lzjson.Type {
	return n.node.Type()
}

// IsObject checks whether the node is valid JSON object.
func (n *Node) IsObject() bool {
	return n.Type() == lzjson.TypeObject
}

// IsArray checks whether the node is valid JSON array.
func (n *Node) IsArray() bool {
	return n.Type() == lzjson.TypeArray
}

// IsString checks whether the node is valid JSON string.
func (n *Node) IsString() bool {
	return n.Type() == lzjson.TypeString
}

// IsNumber checks whether the node is valid JSON number.
func (n *Node) IsNumber() bool {
	return n.Type() == lzjson.TypeNumber
}

// IsBool checks whether the node is valid JSON boolean.
func (n *Node) IsBool() bool {
	return n.Type() == lzjson.TypeBool
}

// IsNull checks whether the node is valid JSON null.
func (n *Node) IsNull() bool {
	return n.Type() == lzjson.TypeNull
}

// IsValid checks whether the node is valid JSON value.
func (n *Node) IsValid() bool {
	return n.Type() != lzjson.TypeError
}

// IsEmpty checks whether the node is having empty value.
func (n *Node) IsEmpty() bool {
	if n.IsObject() {
		return len(n.Keys()) == 0
	}

	return n.Len() == 0
}

// Get gets object's inner value.
// Only works with Object value type.
func (n *Node) Get(key string) *Node {
	return beNode(n.node.Get(key))
}

// GetN gets array's inner value.
// Only works with Array value type.
// 0 for the first item.
func (n *Node) GetN(i int) *Node {
	return beNode(n.node.GetN(i))
}

// Unmarshal parses the JSON node data into variable v.
func (n *Node) Unmarshal(v interface{}) error {
	return n.node.Unmarshal(v)
}

// NewNode reads and decodes a JSON from io.Reader then returns a Node of it.
func NewNode(r io.Reader) *Node {
	return beNode(lzjson.Decode(r))
}

func beNode(n lzjson.Node) *Node {
	return &Node{node: n}
}
