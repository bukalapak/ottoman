package json_test

import (
	bjson "encoding/json"
	"testing"

	"github.com/bukalapak/ottoman/encoding/json"
	"github.com/stretchr/testify/assert"
)

type TestPayload struct {
	N json.Number
	I int
}

type testCaseDecoder struct {
	payload                []byte
	shouldErrUnmarshal     bool
	shouldErrMarshal       bool
	isNumberEmptyOrNullStr bool
	strNVal                string
	shouldErrInt           bool
	intNVal                int64
	floatNVal              float64
	shouldErrFloat         bool
}

type testCaseEncoder struct {
	payload          TestPayload
	shouldErrMarshal bool
	strVal           string
}

func runUnmarshalTest(t *testing.T, tc testCaseDecoder, shouldUseOttomanCoder bool) {

	var err error
	var p TestPayload

	if shouldUseOttomanCoder {
		err = json.Unmarshal(tc.payload, &p)
	} else {
		err = bjson.Unmarshal(tc.payload, &p)
	}

	if tc.shouldErrUnmarshal {
		assert.Error(t, err)
		return
	}

	assert.Nil(t, err)

	if tc.isNumberEmptyOrNullStr { // special condition, will be converted to empty string
		assert.Equal(t, "", p.N.String())
	} else {
		assert.Equal(t, string(tc.strNVal), p.N.String())
	}

	actualInt, err := p.N.Int64()
	if tc.shouldErrInt {
		assert.Error(t, err)
	} else {
		assert.Nil(t, err)
		assert.Equal(t, tc.intNVal, actualInt)
	}

	actualFloat, err := p.N.Float64()
	if tc.shouldErrFloat {
		assert.Error(t, err)
	} else {
		assert.Nil(t, err)
		assert.Equal(t, tc.floatNVal, actualFloat)
	}

}

func runMarshalTest(t *testing.T, tc testCaseEncoder, shouldUseOttomanCoder bool) {
	var err error
	var byteArr []byte

	if shouldUseOttomanCoder {
		byteArr, err = json.Marshal(tc.payload)
	} else {
		byteArr, err = bjson.Marshal(tc.payload)
	}

	if tc.shouldErrMarshal {
		assert.Error(t, err)
		return
	}

	assert.Nil(t, err)

	assert.Equal(t, tc.strVal, string(byteArr))
}

func runUnmarshalTestCases(t *testing.T, tcs []testCaseDecoder, shouldUseOttomanCoder bool) {
	for _, test := range tcs {
		runUnmarshalTest(t, test, shouldUseOttomanCoder)
	}
}

func runMarshalTestCases(t *testing.T, tcs []testCaseEncoder, shouldUseOttomanCoder bool) {
	for _, test := range tcs {
		runMarshalTest(t, test, shouldUseOttomanCoder)
	}
}

func TestNumber_default_decoder(t *testing.T) {
	var testCases = []testCaseDecoder{
		// invalid test
		{payload: []byte(`{"N":"lorem","I":1}`), shouldErrUnmarshal: true},
		{payload: []byte(`{"N":"1L","I":2}`), shouldErrUnmarshal: true},
		{payload: []byte(`{"N":"12.1012.01","I":3}`), shouldErrUnmarshal: true},
		// // empty string
		{payload: []byte(`{"N":"","I":1}`), shouldErrUnmarshal: false, isNumberEmptyOrNullStr: true, strNVal: "", shouldErrInt: true, shouldErrFloat: true},
		{payload: []byte(`{"N":null,"I":2}`), shouldErrUnmarshal: false, isNumberEmptyOrNullStr: true, strNVal: "", shouldErrInt: true, shouldErrFloat: true},
		// valid number
		{payload: []byte(`{"N":"3","I":1}`), shouldErrUnmarshal: false, strNVal: "3", shouldErrInt: false, intNVal: 3, shouldErrFloat: false, floatNVal: 3},
		{payload: []byte(`{"N":"1","I":2}`), shouldErrUnmarshal: false, strNVal: "1", shouldErrInt: false, intNVal: 1, shouldErrFloat: false, floatNVal: 1},
		{payload: []byte(`{"N":"3.6","I":1}`), shouldErrUnmarshal: false, strNVal: "3.6", shouldErrInt: true, shouldErrFloat: false, floatNVal: 3.6},
	}

	runUnmarshalTestCases(t, testCases, false)

}

func TestNumber_default_encoder(t *testing.T) {
	var emptyNumber json.Number

	var testCases = []testCaseEncoder{
		{payload: TestPayload{N: "", I: 1}, shouldErrMarshal: false, strVal: `{"N":"","I":1}`},
		{payload: TestPayload{N: emptyNumber, I: 2}, shouldErrMarshal: false, strVal: `{"N":"","I":2}`},
		{payload: TestPayload{N: "3", I: 2}, shouldErrMarshal: false, strVal: `{"N":"3","I":2}`},
		{payload: TestPayload{N: "3.12", I: 3}, shouldErrMarshal: false, strVal: `{"N":"3.12","I":3}`},
	}

	runMarshalTestCases(t, testCases, false)

}

func TestNumber_ottoman_decoder(t *testing.T) {
	testCases := []testCaseDecoder{
		// // invalid test
		{payload: []byte(`{"N":"lorem","I":1}`), shouldErrUnmarshal: true},
		{payload: []byte(`{"N":"1L","I":2}`), shouldErrUnmarshal: true},
		{payload: []byte(`{"N":"12.1012.01","I":3}`), shouldErrUnmarshal: true},
		// // empty string
		{payload: []byte(`{"N":"","I":1}`), shouldErrUnmarshal: false, isNumberEmptyOrNullStr: true, strNVal: "", shouldErrInt: true, shouldErrFloat: true},
		{payload: []byte(`{"N":null,"I":2}`), shouldErrUnmarshal: false, isNumberEmptyOrNullStr: true, strNVal: "", shouldErrInt: true, shouldErrFloat: true},
		// valid number
		{payload: []byte(`{"N":"3","I":1}`), shouldErrUnmarshal: false, strNVal: "3", shouldErrInt: false, intNVal: 3, shouldErrFloat: false, floatNVal: 3},
		{payload: []byte(`{"N":"1","I":2}`), shouldErrUnmarshal: false, strNVal: "1", shouldErrInt: false, intNVal: 1, shouldErrFloat: false, floatNVal: 1},
		{payload: []byte(`{"N":"3.6","I":1}`), shouldErrUnmarshal: false, strNVal: "3.6", shouldErrInt: true, shouldErrFloat: false, floatNVal: 3.6},
	}

	runUnmarshalTestCases(t, testCases, true) // with ottoman decoder
}

func TestNumber_ottoman_encoder(t *testing.T) {
	var emptyNumber json.Number
	testCases := []testCaseEncoder{
		{payload: TestPayload{N: "", I: 1}, shouldErrMarshal: false, strVal: "{\"N\":\"\",\"I\":1}\n"},
		{payload: TestPayload{N: emptyNumber, I: 2}, shouldErrMarshal: false, strVal: "{\"N\":\"\",\"I\":2}\n"},
		{payload: TestPayload{N: "3", I: 2}, shouldErrMarshal: false, strVal: "{\"N\":\"3\",\"I\":2}\n"},
		{payload: TestPayload{N: "3.12", I: 3}, shouldErrMarshal: false, strVal: "{\"N\":\"3.12\",\"I\":3}\n"},
	}

	runMarshalTestCases(t, testCases, true) // with ottoman decoder
}
