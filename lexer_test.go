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

func TestLex(t *testing.T) {
	input, err := readCAFEFile("./test_data/test-lexer.cafe")
	assert.NoError(t, err)

	lx := newLexer(input)
	lx.lexInput(false)

	expectedNames := []string{
		"// This is a comment",
		"str", `"string"`,
		"// Inline comment",
		"multistr", `"multi" \			"line" \			"string"`,
		"number1", "2023",
		"number2", "3.14159",
		"trickNumber1", `"2023"`,
		"trickNumber2", `"3.14159"`,
		"// These", "// are", "// multiple", "// comments",
		"boolean1", "true",
		"boolean2", "false",
		"trickBoolean1", `"true"`,
		"trickBoolean2", `"false"`,
		"array1", "", // Array open
		`"just"`, `"strings"`, "", // Array close
		"array2", "", // Array open
		`"strings"`, "10", "1.0001", `"and"`, "false", `"other"`, `"types"`, "", // Array close
		"array3", "", // Array open
		`"multi"`, "123", `"line"`, "false", `"array"`, "", // Array close
		"trickArray", `"[this is to trick an array]"`,
		"block1", "", // Block end
		"block2",
		"blockString", `"block string"`,
		"blockInt", "1234",
		"blockFloat", "10.10",
		"blockBool", "false",
		"blockArray", "", // Array start
		`"this"`, "10", `"is"`, "false", `"an"`, `"array"`, `"inside"`, "false", `"of"`, `"a"`, `"block"`,
		"",           // Array end
		"",           // Block end
		"block3", "", // Block end
		"// Comment",
		"block4",
		"nested1",
		"blockNestedString", `"block string"`,
		"blockNestedInt", "1234",
		"blockNestedFloat", "10.10",
		"blockNestedBool", "false",
		"blockNestedArray", "", // Array start
		`"this"`, "10", `"is"`, "false", `"an"`, `"array"`, `"inside"`, "false", `"of"`, `"a"`, `"block"`,
		"", // Array end
		"", // Block end
		"", // Block end
		"arith1", "10 + 10",
		"arith2", "10-8",
		"arith3", "1 * 4",
		"arith4", "10/ 10",
		"compare1", "10 == 20",
		"compare2", "10 != 20",
		"compare3", "10 > 20",
		"compare4", "10 >= 20",
		"compare5", "10 < 20",
		"compare6", "10 <= 20",
		"compare7", "true == true",
		"compare8", "true != true",
		"condition1", "if compare1 : true ? false",
		"// Comment",
		"condition2", "for array1 ...",
	}
	for i, ev := range expectedNames {
		assert.EqualValues(t, ev, lx.items[i].value)
	}
}

func TestLexFunctions(t *testing.T) {
	input, err := readCAFEFile("./test_data/test-functions.cafe")
	assert.NoError(t, err)

	lx := newLexer(input)
	lx.lexInput(false)

	expectedNames := []string{
		"// Strings",
		"testFuncString1", `upper("test function")`,
		"testFuncString2", `lower("TEST FUNCTION")`,
		"testFuncString3", `append("test", " function")`,
		"arr1", "", // Array open
		`"test"`, `"string"`, `"concat"`, `"function"`, "", // Array close
		"testFuncString4", `concat(arr1, " ")`,
		"testFuncString5", `length("test function")`,
		"// Numerical",
		"testFuncNumerical1", "power(5, 2)",
		"testFuncNumerical2", "floor(25, 7)",
		"testFuncNumerical3", "remainder(10, 3)",
		"// Gate Logic",
		"testFuncGateLogic1", "and(true, true)",
		"testFuncGateLogic2", "or(true, false)",
		"testFuncGateLogic3", "nand(true, true)",
		"testFuncGateLogic4", "nor(true, true)",
		"testFuncGateLogic5", "xor(true, false)",
		"testFuncGateLogic6", "xnor(true, true)",
	}
	for i, ev := range expectedNames {
		assert.EqualValues(t, ev, lx.items[i].value)
	}
}
