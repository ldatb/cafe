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
	"fmt"
	"math"
	"strconv"
	"strings"
)

// String functions
func stringFunctions(funcName string, funcParams []string) interface{} {
	// Trim string quotes
	stringValues := make([]string, len(funcParams))
	for i, v := range funcParams {
		stringValues[i] = strings.Trim(strings.TrimSpace(v), `"`)
	}

	switch funcName {
	case "upper":
		return strings.ToUpper(stringValues[0])
	case "lower":
		return strings.ToLower(stringValues[0])
	case "append":
		return stringValues[0] + stringValues[1]
	case "concat":
		// Remove open and close brackets
		// TO-DO
		return nil
		// stringValues[0] = strings.TrimLeft(stringValues[0], "[")
		// stringValues[0] = strings.TrimRight(stringValues[0], "]")
		// array := strings.Split(stringValues[0], ", ")[1:]
		// return strings.Join(array, stringValues[1])
	case "contains":
		return strings.Contains(stringValues[0], stringValues[1])
	case "length":
		return len(stringValues[0])
	default:
		p := fmt.Sprintf("ERROR in parser: function %s is not implemented", funcName)
		panic(p)
	}
}

// Numerical functions
func numericalFunctions(funcName string, funcParams []string) interface{} {
	// Transform parameters into float
	floatParams := make([]float64, len(funcParams))
	hasInt := false
	for i, v := range funcParams {
		// Int
		_, err := strconv.Atoi(strings.TrimSpace(v))
		if err == nil {
			hasInt = true
		}

		// Float
		valFloat, err := strconv.ParseFloat(strings.TrimSpace(v), 32)
		if err == nil {
			precision := math.Pow(10, float64((len(strings.TrimSpace(v)) - 1)))
			floatParams[i] = math.Round(valFloat*precision) / precision
			continue
		}
	}

	switch funcName {
	case "power":
		result := math.Pow(floatParams[0], floatParams[1])
		if hasInt {
			return int(result)
		}
		return result
	case "floor":
		result := math.Floor(floatParams[0] / floatParams[1])
		if hasInt {
			return int(result)
		}
		return result
	case "remainder":
		result := math.Mod(floatParams[0], floatParams[1])
		if hasInt {
			return int(result)
		}
		return result
	default:
		p := fmt.Sprintf("ERROR in parser: function %s is not implemented", funcName)
		panic(p)
	}
}

// Gate logic functions
func gateLogicFunctions(funcName string, funcParams []string) interface{} {
	// Transform parameters into boolean
	boolParams := make([]bool, len(funcParams))
	for i, v := range funcParams {
		valBool, err := strconv.ParseBool(strings.TrimSpace(v))
		if err != nil {
			p := fmt.Sprintf("ERROR in parser: parameter '%s' in function '%s' is not a boolean", v, funcName)
			panic(p)
		}
		boolParams[i] = valBool
	}

	switch funcName {
	case "and":
		return boolParams[0] == boolParams[1]
	case "or":
		return boolParams[0] || boolParams[1]
	case "nand":
		return !(boolParams[0] == boolParams[1])
	case "nor":
		return !(boolParams[0] || boolParams[1])
	case "xor":
		return boolParams[0] != boolParams[1]
	case "xnor":
		return !(boolParams[0] != boolParams[1])
	default:
		p := fmt.Sprintf("ERROR in parser: function %s is not implemented", funcName)
		panic(p)
	}
}
