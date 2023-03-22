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

// Transforms an item with keyMultiString kind
func transformItemMultiString(item string) interface{} {
	// Create an array for the items
	multiStrItems := []string{}

	// Iterate on the item value and remove all "\t" and "\n"
	foundValidElement := false
	validElemStartIndex := 0
	for i, char := range item {
		if foundValidElement && (string(char) == "\t" || string(char) == "\n") {
			multiStrItems = append(multiStrItems, item[validElemStartIndex:i-2])
			foundValidElement = false
		} else if foundValidElement && i == len(item)-1 {
			multiStrItems = append(multiStrItems, item[validElemStartIndex:i])
			foundValidElement = false
		} else if !foundValidElement && (string(char) != "\t" && string(char) != "\n") {
			foundValidElement = true
			validElemStartIndex = i
		}
	}

	// For each element in the array, remove the quote signs
	for i, str := range multiStrItems {
		multiStrItems[i] = strings.Trim(str, `"`)
	}

	// Join the array and return
	return strings.Join(multiStrItems, " ")
}

// Transforms an item with keyArrayStart or keyArrayElem kind
func transformItemArray(item string) interface{} {
	// Remove open and close brackets
	item = strings.TrimLeft(item, "[")
	item = strings.TrimRight(item, "]")

	// Separate elements by comma
	freeElems := strings.Split(item, ", ")[1:]

	// Transform int, float and bool elements
	arrayElems := make([]interface{}, len(freeElems))
	for i, elem := range freeElems {
		// Int
		valInt, err := strconv.Atoi(elem)
		if err == nil {
			arrayElems[i] = valInt
			continue
		}

		// Float
		valFloat, err := strconv.ParseFloat(elem, 32)
		if err == nil {
			precision := math.Pow(10, float64((len(elem) - 1)))
			arrayElems[i] = math.Round(valFloat*precision) / precision
			continue
		}

		// Boolean
		valBool, err := strconv.ParseBool(elem)
		if err == nil {
			arrayElems[i] = valBool
			continue
		}

		// If none, item is a string
		arrayElems[i] = strings.Trim(elem, `"`)
	}
	return arrayElems
}

// Transforms an item with keyArithmetic kind
func transformItemArithmetic(item string) interface{} {
	// Separate all arithmetic symbols and values into an array
	arithmeticArrayOriginal := []interface{}{}
	hasFloat := false
	skipItem := 0
	for i, v := range item {
		// Skip items until all values in the number were passed
		if skipItem != 0 {
			skipItem -= 1
			continue
		}

		// Item is whitespace
		if string(v) == " " {
			continue
		}

		// Arithmetic symbols
		if equalsToMany(string(v), []string{"+", "-", "*", "/"}) {
			// Check if previous item is a number
			if len(arithmeticArrayOriginal) == 0 || equalsToMany(fmt.Sprint(arithmeticArrayOriginal[0]), []string{"+", "-", "*", "/"}) {
				p := fmt.Sprintf("ERROR in parser: arithmetic operation is missing values: %s", item)
				panic(p)
			}

			// Append symbol to arraym
			arithmeticArrayOriginal = append(arithmeticArrayOriginal, string(v))
			continue
		}

		// Numbers
		// Get entire number value
		// First element
		addItemNumber := ""
		if len(arithmeticArrayOriginal) == 0 {
			endOfElem := 0
			for j, c := range item {
				if j > i {
					if equalsToMany(string(c), []string{" ", "+", "-", "*", "/"}) {
						if endOfElem == 0 {
							endOfElem = j
						}
					}
				}
			}
			addItemNumber = item[i:endOfElem]
			skipItem = len(item[i:endOfElem]) - 1
		} else {
			// Last element, so get the entire value until the end of the string
			// If the arithmeticArrayOriginal already has the 3 elements, this will be ignored
			addItemNumber = item[i:]
			skipItem = len(item[i:])
		}

		// Parse number string to float
		var numberVal interface{}
		valInt, err := strconv.Atoi(addItemNumber)
		if err != nil {
			valFloat, err := strconv.ParseFloat(addItemNumber, 32)
			if err != nil {
				p := fmt.Sprintf("ERROR in parser: value in arithmetic operation is not a number: %s", addItemNumber)
				panic(p)
			}
			numberVal = valFloat
			hasFloat = true
		} else {
			numberVal = valInt
		}

		// Check if last element is an arithmetic symbol
		if len(arithmeticArrayOriginal) >= 2 {
			index := len(arithmeticArrayOriginal) - 1
			if !equalsToMany(fmt.Sprint(arithmeticArrayOriginal[index]), []string{"+", "-", "*", "/"}) {
				p := fmt.Sprintf("ERROR in parser: arithmetic operation value must be preceeded by operation symbol: %s", item)
				panic(p)
			}
		}

		// Append number to array
		arithmeticArrayOriginal = append(arithmeticArrayOriginal, numberVal)
	}

	// Perform calculations
	// All calculations are done in float32, even if they're all integers
	// By the end, it will return either as int or float, depending on hasFloat

	// Search for multiplication and division first
	arithmeticArrayFirstOperations := make([]interface{}, len(arithmeticArrayOriginal))
	for i, v := range arithmeticArrayOriginal {
		arithmeticArrayFirstOperations[i] = fmt.Sprint(v)
	}

	for i, v := range arithmeticArrayOriginal {
		if equalsToMany(fmt.Sprint(v), []string{"*", "/"}) {
			// Get values
			val1String := fmt.Sprint(arithmeticArrayOriginal[i-1])
			val2String := fmt.Sprint(arithmeticArrayOriginal[i+1])
			val1, _ := strconv.ParseFloat(val1String, 32)
			val2, _ := strconv.ParseFloat(val2String, 32)

			// Do operation
			var result float64
			if v == "*" {
				result = val1 * val2
			} else if v == "/" {
				result = val1 / val2
			}

			// Transfer result to the array
			arithmeticArrayFirstOperations = arithmeticArrayFirstOperations[i+1:]
			arithmeticArrayFirstOperations[0] = result
		}
	}

	// Do addition and subtraction
	arithmeticArrayOperations := make([]interface{}, len(arithmeticArrayFirstOperations))
	_ = copy(arithmeticArrayOperations, arithmeticArrayFirstOperations)
	for i, v := range arithmeticArrayFirstOperations {
		if equalsToMany(fmt.Sprint(v), []string{"+", "-"}) {
			// Get values
			val1String := fmt.Sprint(arithmeticArrayFirstOperations[i-1])
			val2String := fmt.Sprint(arithmeticArrayFirstOperations[i+1])
			val1, _ := strconv.ParseFloat(val1String, 32)
			val2, _ := strconv.ParseFloat(val2String, 32)

			// Do operation
			var result float64
			if v == "+" {
				result = val1 + val2
			} else if v == "-" {
				result = val1 - val2
			}

			// Transfer result to the array
			arithmeticArrayOperations = arithmeticArrayOperations[i+1:]
			arithmeticArrayOperations[0] = result
		}
	}

	// Get final value of the operation
	if len(arithmeticArrayOperations) > 1 {
		p := fmt.Sprintf("ERROR in parser: extra values in arithmetic operation: %s", item)
		panic(p)
	}
	var result interface{}
	if !hasFloat {
		result, _ = strconv.Atoi(fmt.Sprint(arithmeticArrayOperations[0]))
	} else {
		result, _ = strconv.ParseFloat(fmt.Sprint(arithmeticArrayOperations[0]), 32)
	}

	return result
}

// Transforms an item with keyComparison kind
func transformItemComparison(item string) interface{} {
	// Separate values and comparator to an array
	comparisonArray := []interface{}{}
	skipItem := false
	for i, v := range item {
		if string(v) == " " {
			continue
		}

		// Skip item means that the last element was a comparator with size 2
		if skipItem {
			skipItem = false
			continue
		}

		// Comparator
		if equalsToMany(string(v), []string{"=", "!", ">", "<"}) {
			// Check if it's ==, >= or <=
			itemRange := 1
			if string(item[i+1]) == "=" {
				itemRange += 1
				skipItem = true
			}

			// Check if the array already has an element
			if len(comparisonArray) == 0 {
				p := fmt.Sprintf("ERROR in parser: comparison operation is missing values: %s", item)
				panic(p)
			}

			// Add element to array
			comparisonArray = append(comparisonArray, item[i:i+itemRange])
		} else { // Comparison value
			// If that's the first element, find the next element (comparison symbol or whitespace)
			// to separate it
			if len(comparisonArray) == 0 {
				endOfElem := 0
				for j, c := range item {
					if j > i {
						if equalsToMany(string(c), []string{" ", "=", "!", ">", "<"}) {
							if endOfElem == 0 {
								endOfElem = j
							}
						}
					}
				}
				comparisonArray = append(comparisonArray, item[i:endOfElem])
			} else if len(comparisonArray) > 1 && len(comparisonArray) < 3 {
				// Last element, so get the entire value until the end of the string
				// If the comparisonArray already has the 3 elements, this will be ignored
				comparisonArray = append(comparisonArray, item[i:])
			}
		}
	}

	// A comparison can only have 3 elements, the 2 values and the comparator
	if len(comparisonArray) > 3 {
		p := fmt.Sprintf("ERROR in parser: comparison attributes can only compare 2 items: %s", item)
		panic(p)
	}

	// Check if any of the values is a boolean, if it is, the comparator has to be
	// either == or !=. If it's not the case, panic
	val1Bool, checkVal1Bool := strconv.ParseBool(fmt.Sprint(comparisonArray[0]))
	val2Bool, checkVal2Bool := strconv.ParseBool(fmt.Sprint(comparisonArray[2]))
	if checkVal1Bool == nil || checkVal2Bool == nil {
		if !equalsToMany(fmt.Sprint(comparisonArray[1]), []string{"==", "!="}) {
			p := fmt.Sprintf("ERROR in parser: booleans cannot be compared by %s symbol", fmt.Sprint(comparisonArray[1]))
			panic(p)
		}

		// Both elements have to be a boolean
		if checkVal1Bool != nil || checkVal2Bool != nil {
			p := fmt.Sprintf("ERROR in parser: cannot compare boolean value to numerical value: %s", item)
			panic(p)
		}

		// Compare and return
		switch fmt.Sprint(comparisonArray[1]) {
		case "==":
			return val1Bool == val2Bool
		case "!=":
			return val1Bool != val2Bool
		}
	}

	// Compare numerical values
	val1Float, checkVal1Float := strconv.ParseFloat(fmt.Sprint(comparisonArray[0]), 32)
	val2Float, checkVal2Float := strconv.ParseFloat(fmt.Sprint(comparisonArray[2]), 32)
	if checkVal1Float != nil || checkVal2Float != nil {
		p := fmt.Sprintf("ERROR in parser: can only compare boolean or numerical values: %s", item)
		panic(p)
	}

	switch fmt.Sprint(comparisonArray[1]) {
	case "==":
		return val1Float == val2Float
	case "!=":
		return val1Float != val2Float
	case ">":
		return val1Float > val2Float
	case ">=":
		return val1Float >= val2Float
	case "<":
		return val1Float < val2Float
	case "<=":
		return val1Float <= val2Float
	}
	return false
}

// Transforms an item with keyFunction kind
func transformItemFunction(item string) interface{} {
	// Get function name
	funcNameIndex := strings.Index(item, "(")
	funcName := strings.TrimSpace(item[:funcNameIndex])

	// Get first attribute
	funcParams := []string{}

	// Single parameter functions
	singleParamFunctions := []string{"upper", "lower", "length"}
	if equalsToMany(funcName, singleParamFunctions) {
		funcParams = append(funcParams, item[funcNameIndex+1:len(item)-1])
	} else { // Multiple parameter functions
		for i, v := range item {
			if i < funcNameIndex {
				continue
			}

			// Ignore comma if it's inside brackets
			insideBrackets := false
			if string(v) == `"` && !insideBrackets {
				insideBrackets = true
			} else if string(v) == `"` && insideBrackets {
				insideBrackets = false
			}

			// Check for comma
			if string(v) == "," && !insideBrackets {
				funcParams = append(funcParams, item[funcNameIndex+1:i])
				funcParams = append(funcParams, item[i+1:len(item)-1])
				break
			}
		}
	}

	// Strings
	stringFunctionNames := []string{"upper", "lower", "append", "concat", "contains", "length"}
	if equalsToMany(funcName, stringFunctionNames) {
		return stringFunctions(funcName, funcParams)
	}

	// Numerical
	numericalFunctionNames := []string{"power", "floor", "remainder"}
	if equalsToMany(funcName, numericalFunctionNames) {
		return numericalFunctions(funcName, funcParams)
	}

	// Gate logic
	gateLogicFunctionNames := []string{"and", "or", "nand", "nor", "xor", "xnor"}
	if equalsToMany(funcName, gateLogicFunctionNames) {
		return gateLogicFunctions(funcName, funcParams)
	}

	// Panic
	p := fmt.Sprintf("ERROR in parser: function %s is not implemented", funcName)
	panic(p)
}

// Transforms an item's value string into an interface
func transformItem(item string, kind keyKind) interface{} {
	// String
	if kind == keyString {
		// Simply remove the quote signs and return
		return strings.Trim(item, `"`)
	}

	// Multiline string
	if kind == keyMultiString {
		return transformItemMultiString(item)
	}

	// Int
	if kind == keyInt {
		val, err := strconv.Atoi(item)
		if err != nil {
			p := fmt.Sprintf("ERROR in parser: non int item %s tried to be parsed as int", item)
			panic(p)
		}
		return val
	}

	// Float
	if kind == keyFloat {
		val, err := strconv.ParseFloat(item, 32)
		if err != nil {
			p := fmt.Sprintf("ERROR in parser: non float item %s tried to be parsed as float", item)
			panic(p)
		}
		precision := math.Pow(10, float64((len(item) - 1)))
		return math.Round(val*precision) / precision
	}

	// Bool
	if kind == keyBool {
		val, err := strconv.ParseBool(item)
		if err != nil {

			p := fmt.Sprintf("ERROR in parser: non boolean item %s tried to be parsed as boolean", item)
			panic(p)
		}
		return val
	}

	// Array
	if kind == keyArrayStart || kind == keyArrayElem {
		return transformItemArray(item)
	}

	// Arithmetic
	if kind == keyArithmetic {
		return transformItemArithmetic(item)
	}

	// Comparison
	if kind == keyComparison {
		return transformItemComparison(item)
	}

	// Condition
	if kind == keyCondition {
		// TO-DO
		return nil
	}

	// Function
	if kind == keyFunction {
		return transformItemFunction(item)
	}

	// Unknown
	p := fmt.Sprintf("ERROR in parser: unknown item kind: %s", keyKindStr(kind))
	panic(p)
}
