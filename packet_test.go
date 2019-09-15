package fh4server

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type parseTestCase struct {
	// base64 encoded packet bytes (from logs)
	input          string
	expectedFields map[string]interface{}
	expectedTags   map[string]interface{}
}

func TestParseForValidCases(t *testing.T) {
	testCases := []parseTestCase{
		parseTestCase{
			input: "AQAAAHGpEAD2//lF+P9HREg17ETOUljBxAq0vp6jGk" +
				"AHMB9BgzqMPgAOVj/faPM9fXSzv8Yygry3H/w/9+gKvWx7s7vyegc/NSD9P" +
				"lpD7D7ZpEo/43IfQG+03z/mgKtAzkzePwZWCz8L5qFASKvOQXFkYUEAAAAA" +
				"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAJmZGT+ZmRk/mZkZP5mZGT8" +
				"b06vAHRbewF+x7MBXmtjAWWu9QBYF5UDkJRJBbp7fQNDBPjzQQQg8IFO7O6" +
				"QsFz0kBAAAAwAAACADAAACAAAABAAAACMAAAAAAAAAAAAAAJeqw8Lv3WtDx" +
				"9qxxRPPH0EAAAAAAAAAAGVnm0LYSZpC/FaeQvxWnkLMdrLAAACAPwAAAAAA" +
				"AAAAAAAAAAAAAABs2xdEAAAA/wAAAAGBAAAA",
			expectedFields: map[string]interface{}{
				"accel":      255,
				"brake":      0,
				"clutch":     0,
				"gear":       1,
				"hand_brake": 0,
				"speed":      "9.988055",
			},
			expectedTags: map[string]interface{}{
				"is_race_on":           "1",
				"car_class":            "3",
				"num_engine_cylinders": "4",
			},
		},
	}
	r := require.New(t)
	for testCaseIndex, testCase := range testCases {
		packetBytes, err := base64.StdEncoding.DecodeString(testCase.input)
		if err != nil {
			t.Errorf("invalid input for test case %d: %v", testCaseIndex, err)
		}
		actual := Parse(packetBytes, func(label string) bool {
			// Use all of the fields specified in the expected val as keys in
			// the whitelist function.
			_, aField := testCase.expectedFields[label]
			_, aTag := testCase.expectedTags[label]
			return aField || aTag
		})

		checkMap(r, testCase.expectedFields, actual.Fields)

		actualTags := make(map[string]interface{})
		for k, v := range actual.Tags {
			actualTags[k] = v
		}
		checkMap(r, testCase.expectedTags, actualTags)
	}
}

func checkMap(r *require.Assertions, expected, actual map[string]interface{}) {
	for k := range expected {
		expectedStr := fmt.Sprintf("%v", expected[k])
		actualVal, _ := actual[k]
		r.Contains(actual, k, "missing key")

		actualStr := fmt.Sprintf("%v", actualVal)
		r.Equal(expectedStr, actualStr)
	}
}
