package json_test

import (
	bjson "encoding/json"
	"strings"
	"testing"

	"github.com/bukalapak/ottoman/encoding/json"
	"github.com/stretchr/testify/assert"
)

var sample = `{"numberK":123,"numberStringK":"123","emptyStringK":"","nullK":null}`
var expectedString = `{"numberK":123,"numberStringK":123,"emptyStringK":0,"nullK":0}`

type NumberKind struct {
	Number       bjson.Number `json:"numberK"`
	NumberString bjson.Number `json:"numberStringK"`
	EmptyString  bjson.Number `json:"emptyStringK"`
	Null         bjson.Number `json:"nullK"`
}

func TestOriginalJSONNumber(t *testing.T) {
	x := NumberKind{
		Null:         bjson.Number(""),
		Number:       bjson.Number("123"),
		NumberString: bjson.Number("123"),
		EmptyString:  bjson.Number(""),
	}

	var v NumberKind

	dec := bjson.NewDecoder(strings.NewReader(sample))
	dec.UseNumber()

	err := dec.Decode(&v)
	assert.Contains(t, err.Error(), "json: invalid number literal")
	assert.Equal(t, x, v)

	ss, err := json.Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, expectedString, strings.TrimSpace(string(ss)))
}

type CNumberKind struct {
	Number       json.Number `json:"numberK"`
	NumberString json.Number `json:"numberStringK"`
	EmptyString  json.Number `json:"emptyStringK"`
	Null         json.Number `json:"nullK"`
}

func TestCustomJSONNumber(t *testing.T) {
	x := CNumberKind{
		Null:         json.Number(""),
		Number:       json.Number("123"),
		NumberString: json.Number("123"),
		EmptyString:  json.Number(""),
	}

	var v CNumberKind

	dec := json.NewDecoder(strings.NewReader(sample))
	err := dec.Decode(&v)
	assert.Nil(t, err)
	assert.Equal(t, x, v)

	ss, err := json.Marshal(&v)
	assert.Nil(t, err)
	assert.Equal(t, expectedString, strings.TrimSpace(string(ss)))
}

type StringKind struct {
	Original bjson.Number `json:"original"`
	Custom   json.Number  `json:"custom"`
}

func TestStringKind(t *testing.T) {
	var v StringKind

	sO := `{"Original":"abc","custom":"123"}`
	decO := bjson.NewDecoder(strings.NewReader(sO))
	err := decO.Decode(&v)
	assert.Contains(t, err.Error(), "json: invalid number literal")

	sC := `{"Original":"123","custom":"abc"}`
	decC := json.NewDecoder(strings.NewReader(sC))
	err = decC.Decode(&v)
	assert.Contains(t, err.Error(), "json: invalid number literal")
}
