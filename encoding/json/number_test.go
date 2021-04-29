package json_test

import (
	"testing"

	bjson "encoding/json"

	"github.com/bukalapak/ottoman/encoding/json"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	test               string
	shouldErrUnmarshal bool
	strVal             string
	shouldErrInt       bool
	intVal             int64
	floatVal           float64
	shouldErrFloat     bool
}

func runTest(t *testing.T, tc testCase, shouldUseOttomanCoder bool) {

	var err error
	var n json.Number

	if shouldUseOttomanCoder {
		err = json.Unmarshal([]byte(tc.test), &n)
	} else {
		err = n.UnmarshalJSON([]byte(tc.test))
	}

	if tc.shouldErrUnmarshal {
		assert.Error(t, err)
		return
	}

	var errM error
	var byteArr []byte
	if shouldUseOttomanCoder {
		byteArr, errM = json.Marshal(n)
	} else {
		byteArr, errM = bjson.Marshal(n)
	}

	assert.Nil(t, errM)

	if tc.test == "" || tc.test == "null" { // special condition, will be converted to empty string
		assert.Equal(t, "", string(byteArr))
	} else {
		assert.Equal(t, tc.test, string(byteArr))
	}

	assert.Equal(t, tc.strVal, n.Number.String())

	actualInt, err := n.Number.Int64()
	if tc.shouldErrInt {
		assert.Error(t, err)
	} else {
		assert.Nil(t, err)
		assert.Equal(t, tc.intVal, actualInt)
	}

	actualFloat, err := n.Number.Float64()
	if tc.shouldErrFloat {
		assert.Error(t, err)
	} else {
		assert.Nil(t, err)
		assert.Equal(t, tc.floatVal, actualFloat)
	}

}

func runTestCases(t *testing.T, tcs []testCase, shouldUseOttomanCoder bool) {
	for _, test := range tcs {
		runTest(t, test, shouldUseOttomanCoder)
	}
}

func TestNumber_invalid(t *testing.T) {

	var testCases = []testCase{
		{test: "lorem", shouldErrUnmarshal: true},
		{test: "1L", shouldErrUnmarshal: true},
		{test: "12.1012.01", shouldErrUnmarshal: true},
	}

	runTestCases(t, testCases, false)
	runTestCases(t, testCases, true) // with ottoman decoder

}

func TestNumber_empty_string(t *testing.T) {

	var testCases = []testCase{
		{test: "", shouldErrUnmarshal: false, shouldErrInt: true, shouldErrFloat: true},
		{test: "null", shouldErrUnmarshal: false, shouldErrInt: true, shouldErrFloat: true},
	}

	runTestCases(t, testCases, false)
	runTestCases(t, testCases, true) // with ottoman decoder

}

func TestNumber_number_string(t *testing.T) {

	var testCases = []testCase{
		{test: "3", shouldErrUnmarshal: false, strVal: "3", shouldErrInt: false, intVal: 3, shouldErrFloat: false, floatVal: 3},
		{test: "1", shouldErrUnmarshal: false, strVal: "1", shouldErrInt: false, intVal: 1, shouldErrFloat: false, floatVal: 1},
		{test: "3.6", shouldErrUnmarshal: false, strVal: "3.6", shouldErrInt: true, shouldErrFloat: false, floatVal: 3.6},
	}

	runTestCases(t, testCases, false)
	runTestCases(t, testCases, true) // with ottoman decoder
}
