// Copyright (C) 2023 Lucas de Ataides <lucasatab@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package cafe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserGlobalAttrributes(t *testing.T) {
	input, err := readCAFEFile("./test_data/test-lexer.cafe")
	assert.NoError(t, err)

	p := newParser(input)
	p.parseItems(false)

	expectedMap := map[string]attribute{
		"str": {
			Name:  "str",
			Value: "string",
			kind:  attrString,
		},
		"multistr": {
			Name:  "multistr",
			Value: "multi line string",
			kind:  attrString,
		},
		"number1": {
			Name:  "number1",
			Value: 2023,
			kind:  attrInt,
		},
		"number2": {
			Name:  "number2",
			Value: 3.14159,
			kind:  attrFloat,
		},
		"trickNumber1": {
			Name:  "trickNumber1",
			Value: "2023",
			kind:  attrString,
		},
		"trickNumber2": {
			Name:  "trickNumber2",
			Value: "3.14159",
			kind:  attrString,
		},
		"boolean1": {
			Name:  "boolean1",
			Value: true,
			kind:  attrBool,
		},
		"boolean2": {
			Name:  "boolean2",
			Value: false,
			kind:  attrBool,
		},
		"trickBoolean1": {
			Name:  "trickBoolean1",
			Value: "true",
			kind:  attrString,
		},
		"trickBoolean2": {
			Name:  "trickBoolean2",
			Value: "false",
			kind:  attrString,
		},
		"array1": {
			Name:  "array1",
			Value: []interface{}{"just", "strings"},
			kind:  attrArray,
		},
		"array2": {
			Name:  "array2",
			Value: []interface{}{"strings", 10, 1.0001, "and", false, "other", "types"},
			kind:  attrArray,
		},
		"trickArray": {
			Name:  "trickArray",
			Value: "[this is to trick an array]",
			kind:  attrString,
		},
		"arith1": {
			Name:  "arith1",
			Value: 20,
			kind:  attrArithmetic,
		},
		"arith2": {
			Name:  "arith2",
			Value: 2,
			kind:  attrArithmetic,
		},
		"arith3": {
			Name:  "arith3",
			Value: 4,
			kind:  attrArithmetic,
		},
		"arith4": {
			Name:  "arith4",
			Value: 1,
			kind:  attrArithmetic,
		},
		"compare1": {
			Name:  "compare1",
			Value: false,
			kind:  attrComparison,
		},
		"compare2": {
			Name:  "compare2",
			Value: true,
			kind:  attrComparison,
		},
		"compare3": {
			Name:  "compare3",
			Value: false,
			kind:  attrComparison,
		},
		"compare4": {
			Name:  "compare4",
			Value: false,
			kind:  attrComparison,
		},
		"compare5": {
			Name:  "compare5",
			Value: true,
			kind:  attrComparison,
		},
		"compare6": {
			Name:  "compare6",
			Value: true,
			kind:  attrComparison,
		},
		"compare7": {
			Name:  "compare7",
			Value: true,
			kind:  attrComparison,
		},
		"compare8": {
			Name:  "compare8",
			Value: false,
			kind:  attrComparison,
		},
	}

	for _, v := range expectedMap {
		assert.Equal(t, v, p.Attributes[v.Name])
	}
}

func TestParserBlockAttributes(t *testing.T) {
	input, err := readCAFEFile("./test_data/test-lexer.cafe")
	assert.NoError(t, err)

	p := newParser(input)
	p.parseItems(false)

	expectedMap := map[string]attribute{
		"blockString": {
			Name:  "blockString",
			Value: "block string",
			kind:  attrString,
		},
		"blockInt": {
			Name:  "blockInt",
			Value: 1234,
			kind:  attrInt,
		},
		"blockFloat": {
			Name:  "blockFloat",
			Value: 10.10,
			kind:  attrFloat,
		},
		"blockBool": {
			Name:  "blockBool",
			Value: false,
			kind:  attrBool,
		},
		"blockArray": {
			Name:  "blockArray",
			Value: []interface{}{"this", 10, "is", false, "an", "array", "inside", false, "of", "a", "block"},
			kind:  attrArray,
		},
	}

	// First block
	for _, v := range expectedMap {
		assert.Equal(t, v, p.Blocks["block2"].Attributes[v.Name])
	}

	expectedNestedMap := map[string]attribute{
		"blockNestedString": {
			Name:  "blockNestedString",
			Value: "block string",
			kind:  attrString,
		},
		"blockNestedInt": {
			Name:  "blockNestedInt",
			Value: 1234,
			kind:  attrInt,
		},
		"blockNestedFloat": {
			Name:  "blockNestedFloat",
			Value: 10.10,
			kind:  attrFloat,
		},
		"blockNestedBool": {
			Name:  "blockNestedBool",
			Value: false,
			kind:  attrBool,
		},
		"blockNestedArray": {
			Name:  "blockNestedArray",
			Value: []interface{}{"this", 10, "is", false, "an", "array", "inside", false, "of", "a", "block"},
			kind:  attrArray,
		},
	}

	// Nested block
	for _, v := range expectedNestedMap {
		assert.Equal(t, v, p.Blocks["block4"].Blocks["nested1"].Attributes[v.Name])
	}
}

func TestParseFunctions(t *testing.T) {
	input, err := readCAFEFile("./test_data/test-functions.cafe")
	assert.NoError(t, err)

	p := newParser(input)
	p.parseItems(false)

	expectedMap := map[string]attribute{
		// String functions
		"testFuncString1": {
			Name:  "testFuncString1",
			Value: "TEST FUNCTION",
			kind:  attrFunction,
		},
		"testFuncString2": {
			Name:  "testFuncString2",
			Value: "test function",
			kind:  attrFunction,
		},
		"testFuncString3": {
			Name:  "testFuncString3",
			Value: "test function",
			kind:  attrFunction,
		},
		// "testFuncString4": {
		// 	Name:  "testFuncString4",
		// 	Value: "test string concat function",
		// 	kind:  attrFunction,
		// },
		"testFuncString5": {
			Name:  "testFuncString5",
			Value: 13,
			kind:  attrFunction,
		},
		// Numeric functions
		"testFuncNumerical1": {
			Name:  "testFuncNumerical1",
			Value: 25,
			kind:  attrFunction,
		},
		"testFuncNumerical2": {
			Name:  "testFuncNumerical2",
			Value: 3,
			kind:  attrFunction,
		},
		"testFuncNumerical3": {
			Name:  "testFuncNumerical3",
			Value: 1,
			kind:  attrFunction,
		},
		// Gate logic functions
		"testFuncGateLogic1": {
			Name:  "testFuncGateLogic1",
			Value: true,
			kind:  attrFunction,
		},
		"testFuncGateLogic2": {
			Name:  "testFuncGateLogic2",
			Value: true,
			kind:  attrFunction,
		},
		"testFuncGateLogic3": {
			Name:  "testFuncGateLogic3",
			Value: false,
			kind:  attrFunction,
		},
		"testFuncGateLogic4": {
			Name:  "testFuncGateLogic4",
			Value: false,
			kind:  attrFunction,
		},
		"testFuncGateLogic5": {
			Name:  "testFuncGateLogic5",
			Value: true,
			kind:  attrFunction,
		},
		"testFuncGateLogic6": {
			Name:  "testFuncGateLogic6",
			Value: true,
			kind:  attrFunction,
		},
	}

	for _, v := range expectedMap {
		assert.Equal(t, v, p.Attributes[v.Name])
	}
}
