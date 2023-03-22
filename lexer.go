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
	"strconv"
	"strings"
)

// keyKind defines all kinds of possible keys in an CAFE file
type keyKind int

const (
	keyNIL         keyKind = iota // 0, no type
	keyEOF                        // 1
	keyTab                        // 2
	keyError                      // 3
	keyComment                    // 4
	keyAttrDef                    // 5
	keyBlockStart                 // 6
	keyBlockEnd                   // 7
	keyAttrCall                   // 8
	keyString                     // 9
	keyMultiString                // 10
	keyInt                        // 11
	keyFloat                      // 12
	keyBool                       // 13
	keyArrayStart                 // 14
	keyArrayEnd                   // 15
	keyArrayElem                  // 16
	keyArithmetic                 // 17
	keyComparison                 // 18
	keyCondition                  // 19
	keyFunction                   // 20
)

// position is the position of the parser.
type position struct {
	// Line number starting in 1
	Line int

	// First byte
	Start int

	// Length of the bytes
	Length int
}

// An item differs to a prototype because it requires the value to
// be an array of strings, instead of the whole string
type prototype struct {
	// Kind of the item
	kind keyKind

	// Item value in array of string
	value []string

	// position of the item
	position position
}

// An item is a collection of bytes in an CAFE file
// It corresponds to a definition, either of an attribute or a block
type item struct {
	// Kind of the item
	kind keyKind

	// Item value in string
	value string

	// position of the item
	position position
}

// Lexer reads an array of bytes and decodes it to an array of items
type lexer struct {
	// The raw input corresponding to an CAFE file
	input []string

	// All items of the CAFE file
	items []item

	// Prototype of the next item
	proto prototype

	// Current line of the lexer
	currentLine int

	// Index of the current byte
	currentByteIndex int

	// Current byte of the lexer
	currentByte string

	// Index of the last EOL
	lastEOL int

	// Emits true when in EOF
	atEOF bool
}

// Creates a lexer
func newLexer(input []string) *lexer {
	return &lexer{
		input:            input,
		items:            []item{},
		proto:            prototype{kind: keyNIL},
		currentLine:      1,
		currentByteIndex: 0,
		currentByte:      input[0],
		lastEOL:          0,
		atEOF:            false,
	}
}

// Gets the start of the next item of the input in a lexer
// Moves the lexer to the next item start
// Returns the first byte of the next item
func (l *lexer) next(whitespace bool, eol bool) {
	// Add the prototype to the items list and reset the prototype
	if l.proto.kind != keyNIL {
		// Add required values to proto's position
		l.proto.position.Line = l.currentLine
		l.proto.position.Start = l.currentByteIndex

		// Create new item based on the prototype
		newItem := item{
			kind:     l.proto.kind,
			value:    strings.TrimSpace(strings.Join(l.proto.value, "")),
			position: l.proto.position,
		}

		l.items = append(l.items, newItem)
		l.proto = prototype{kind: keyNIL}
	}

	lastItem := l.items[len(l.items)-1]
	lastItemLastByte := lastItem.position.Start + lastItem.position.Length
	if lastItemLastByte == len(l.input)+1 {
		panic("ERROR in lexer: next called at EOF")
	}

	// Check EOL
	if eol {
		l.currentLine += 1
		_, lastEOLIndex := l.peekAndFind("\n")
		l.lastEOL = lastEOLIndex
	}

	// Check EOF
	if l.currentByteIndex == len(l.input)-1 {
		l.atEOF = true
	}

	if !l.atEOF {
		if whitespace {
			l.currentByteIndex += 1
		} else {
			l.currentByteIndex = lastItemLastByte + 1
		}

		// Check EOF after advancing the bytes
		if l.currentByteIndex == len(l.input) {
			l.atEOF = true
		} else {
			l.currentByte = l.input[l.currentByteIndex]
		}
	}
}

// Gets the start of the next item of the input in a lexer
// DOES NOT MOVE THE LEXER
// Returns the first byte of the next item
func (l *lexer) peek() string {
	return l.input[l.currentByteIndex+1]
}

// Peek and search for a key until EOL
// Here the argument is a simple string since
// not necessarily a key needs to be found
// If the key was not found, returns (false, 0)
// If the key was found, returns (true, index)
func (l *lexer) peekAndFind(key string) (bool, int) {
	for i, v := range l.input {
		if i > l.currentByteIndex {
			if v == "\n" && key != "\n" {
				return false, 0
			}
			if v == key {
				return true, i
			}
		}
	}
	return false, 0
}

// Peek and search for a key until EOL
// Only start looking after the given index
// If the key was not found, returns (false, 0)
// If the key was found, returns (true, index)
func (l *lexer) peekAndFindAfter(key string, startLookingAfter int) (bool, int) {
	for i, v := range l.input {
		if i > l.currentByteIndex && i > startLookingAfter {
			if v == "\n" && key != "\n" {
				return false, 0
			}
			if v == key {
				return true, i
			}
		}
	}
	return false, 0
}

// Similar to peekAndFind, this will search for multiple keys
// until EOL
// For this function, a check for EOL is not made
func (l *lexer) peekAndFindMany(keys []string) bool {
	for i, v := range l.input {
		if i > l.currentByteIndex {
			if v == "\n" {
				return false
			}
			for _, key := range keys {
				if v == key {
					return true
				}
			}
		}
	}
	return false
}

// Find next EOL index (Unicode U+000A)
// Do not use peekAndFind here as EOL can be next to EOF,
// which would cause errors if peekAndFind was used
func (l *lexer) peekEOL() int {
	for i, v := range l.input {
		if i > l.currentByteIndex {
			if v == "\n" {
				return i
			}
			if i == len(l.input)-1 {
				return i
			}
		}
	}
	return 0
}

// Find next comment before EOL
func (l *lexer) peekComment(startLookingAfter int) (bool, int) {
	for i, v := range l.input {
		if i > l.currentByteIndex && i > startLookingAfter {
			if v == "\n" {
				return false, 0
			}
			if v == "/" {
				if l.input[i+1] == "/" {
					return true, i
				}
			}
		}
	}
	return false, 0
}

// Checks the previous byte on the input
func (l *lexer) previousByte() string {
	if l.currentByteIndex == 0 {
		return "\n"
	}
	return l.input[l.currentByteIndex-1]
}

// Checks the previous item of the byte
// Returns an item
func (l *lexer) previousItem() *item {
	return &l.items[len(l.items)-1]
}

// Searches valid characters between an index and an EOL in the lexer
func (l *lexer) searchNonWhitespacedValue() (string, int, int) {
	firstIndex := 0
	lastIndex := 0
	for i, v := range l.input {
		if i >= l.currentByteIndex {
			// Get first non-whitespace
			if firstIndex == 0 && (v != " " && v != "\n") {
				if firstIndex == 0 {
					firstIndex = i
				}
			}

			// Find last whitespace and set the lastIndex to the
			// previous index
			if (firstIndex != 0 && lastIndex == 0) && (v == " " || v == "\n") {
				lastIndex = i
				break
			}

			// Last rune in the input
			if i == len(l.input)-1 {
				lastIndex = i + 1
				break
			}
		}
	}

	// No value found
	if lastIndex == 0 {
		return "", 0, 0
	}

	// Return concatenate value
	return strings.TrimSpace(strings.Join(l.input[firstIndex:lastIndex], "")), firstIndex, lastIndex
}

// End of line (Unicode U+000A)
func (l *lexer) lexEOL() bool {
	if l.currentByte == "\n" && l.currentByteIndex != 0 {
		l.next(true, true)
		return true
	}
	return false
}

// Tab (Unicode U+0009)
func (l *lexer) lexTab() bool {
	next4Characters := strings.Join(l.input[l.currentByteIndex:l.currentByteIndex+3], "")
	if next4Characters == "    " {
		// Call next 4 times to skip the tab
		for i := 0; i < 4; i++ {
			l.next(true, false)
		}
		return true
	}
	return false
}

// End of file
func (l *lexer) lexEOF() bool {
	if l.currentByteIndex == len(l.input)-1 {
		l.atEOF = true
	}
	return false
}

// Whitespace
func (l *lexer) lexWhitespace() bool {
	if l.currentByte == " " {
		l.next(true, false)
		return true
	}
	return false
}

// Comment definition
// Comments don't have to be preceded by an EOL,
// but must be finished by one
func (l *lexer) lexComment() bool {
	if l.currentByte != "/" {
		return false
	}

	if l.peek() != "/" {
		return false
	}

	// Find EOL
	nextEOL := l.peekEOL()

	// Create item as a prototype and call next
	l.proto = prototype{
		kind:  keyComment,
		value: l.input[l.currentByteIndex:nextEOL],
		position: position{
			Length: nextEOL - l.currentByteIndex,
		},
	}

	l.next(false, true)
	return true
}

// Attribute definition
// It has to preceeded by an EOL or whitespaces only
func (l *lexer) lexAttributeDef() bool {
	// Can't be proceeded by keyAttrDef
	if len(l.items) != 0 {
		if l.previousItem().kind == keyAttrDef {
			return false
		}
	}

	// Has to be proceeded by EOL or whitespaces
	if l.previousByte() != "\n" {
		for i, c := range l.input {
			if i > l.lastEOL && i < l.currentByteIndex {
				if c != " " {
					return false
				}
			}
		}
	}

	hasEqual, equalIndex := l.peekAndFind("=")
	if !hasEqual {
		return false
	}

	l.proto = prototype{
		kind:  keyAttrDef,
		value: l.input[l.currentByteIndex:equalIndex],
		position: position{
			Length: equalIndex - l.currentByteIndex,
		},
	}

	l.next(false, false)
	return true
}

// Start of array
// Has to be preceeded by a keyAttrDef
func (l *lexer) lexArrayStart() bool {
	// Has to be proceeded by keyAttrDef
	if l.previousItem().kind != keyAttrDef {
		return false
	}

	// Check if this isn't actually a string
	if l.currentByte == `"` {
		return false
	}

	// Assume current index is a opening bracket ([)
	// If it's not, peekAndFind will look and return if
	// any opening bracket is found
	hasOpenBracket := true
	openBracketIndex := l.currentByteIndex
	if l.currentByte != `[` {
		hasOpenBracket, openBracketIndex = l.peekAndFind(`[`)
	}
	if !hasOpenBracket {
		return false
	}

	// Create prototype and call next
	l.proto = prototype{
		kind:  keyArrayStart,
		value: l.input[l.currentByteIndex:openBracketIndex],
		position: position{
			Length: openBracketIndex - l.currentByteIndex,
		},
	}
	l.next(false, false)
	return true
}

// End of array
// Has to be preceeded by a keyArrayStart (empty array) or
// by a keyArrayElem
func (l *lexer) lexArrayEnd() bool {
	// by a keyArrayElem
	// Has to be proceeded by keyArrayStart or keyArrayElem
	if l.previousItem().kind != keyArrayStart && l.previousItem().kind != keyArrayElem {
		return false
	}

	// Check if this isn't actually a string
	if l.currentByte == `"` {
		return false
	}

	// Check if there's no comma, meaning that there's still elements left in the array
	hasComma, _ := l.peekAndFind(",")
	if hasComma {
		return false
	}

	// Assume current index is a closing bracket (])
	// If it's not, peekAndFind will look and return if
	// any closing bracket is found
	hasCloseBracket := true
	closeIndex := l.currentByteIndex
	if l.currentByte != `]` {
		hasCloseBracket, closeIndex = l.peekAndFind(`]`)
	}
	if !hasCloseBracket {
		return false
	}

	// Check if EOL comes before the closing bracket
	// If that's true, this is probably the last element
	// in a multiline array
	EOLIndex := l.peekEOL()
	if EOLIndex < closeIndex {
		return false
	}

	// Create prototype and call next
	l.proto = prototype{
		kind:  keyArrayEnd,
		value: l.input[l.currentByteIndex:closeIndex],
		position: position{
			Length: closeIndex - l.currentByteIndex,
		},
	}
	l.next(false, false)
	return true
}

// Array element
// Has to be preceeded by a keyArrayStart or keyArrayElem
// and ended by a comma (,) or a closing bracket
func (l *lexer) lexArrayElem() bool {
	// Has to be proceeded by keyArrayStart or keyArrayElem
	if l.previousItem().kind != keyArrayStart && l.previousItem().kind != keyArrayElem {
		return false
	}

	// Search for next comma or end of array
	isEOL := false
	lengthModifier := 0
	hasEndOfElem, endOfArrayElem := l.peekAndFind(",")
	if !hasEndOfElem {
		hasEndOfElem, endOfArrayElem = l.peekAndFind("]")
		lengthModifier = 1
	}

	// If still no end of element was found, this probably is the last
	// element of a multiline array
	if !hasEndOfElem {
		endOfArrayElem = l.peekEOL()
		isEOL = true
	}

	// Create prototype and call next
	l.proto = prototype{
		kind:  keyArrayElem,
		value: l.input[l.currentByteIndex:endOfArrayElem],
		position: position{
			Length: endOfArrayElem - l.currentByteIndex - lengthModifier,
		},
	}
	l.next(false, isEOL)
	return true
}

// ATTRIBUTE TYPES
// These have to be preceeded by a keyAttrDef

// Attribute type: Call to another attribute (keyAttrCall)
func (l *lexer) lexAttrCall() bool {
	// Has to be proceeded by keyAttrDef
	if l.previousItem().kind != keyAttrDef {
		return false
	}

	// Get all characters between whitespace
	checkAttrCall, firstIndex, lastIndex := l.searchNonWhitespacedValue()
	if checkAttrCall == "" {
		return false
	}

	// Create prototype and call next
	l.proto = prototype{
		kind:  keyAttrCall,
		value: l.input[firstIndex:lastIndex],
		position: position{
			Length: lastIndex - l.currentByteIndex,
		},
	}
	l.next(false, false)
	return true
}

// Attribute type: String (keyString)
func (l *lexer) lexAttrString() bool {
	// Has to be proceeded by keyAttrDef
	if l.previousItem().kind != keyAttrDef {
		return false
	}

	// Assume current index is an opening quote (")
	// If it's not, peekAndFind will look and return if
	// any opening quote is found
	hasOpenQuote := true
	openQuoteIndex := l.currentByteIndex
	if l.currentByte != `"` {
		hasOpenQuote, openQuoteQuoteSearch := l.peekAndFind(`"`)
		if hasOpenQuote {
			openQuoteIndex = openQuoteQuoteSearch
		}
	}
	if !hasOpenQuote {
		return false
	}

	// Discover if there's any closing quote between the first quote and EOL
	hasCloseQuote, closeQuoteIndex := l.peekAndFindAfter(`"`, openQuoteIndex)
	if !hasCloseQuote {
		return false
	}

	// Create prototype and call next
	l.proto = prototype{
		kind:  keyString,
		value: l.input[openQuoteIndex : closeQuoteIndex+1],
		position: position{
			Length: closeQuoteIndex - l.currentByteIndex,
		},
	}
	l.next(false, false)
	return true
}

// Attribute type: Multiline string (keyMultiString)
func (l *lexer) lexAttrMultiString() bool {
	// Has to be proceeded by keyAttrDef
	if l.previousItem().kind != keyAttrDef {
		return false
	}

	// Assume current index is an opening quote (")
	// If it's not, peekAndFind will look and return if
	// any opening quote is found
	hasOpenQuote := true
	openQuoteIndex := l.currentByteIndex
	if l.currentByte != `"` {
		hasOpenQuote, openQuoteQuoteSearch := l.peekAndFind(`"`)
		if hasOpenQuote {
			openQuoteIndex = openQuoteQuoteSearch
		}
	}
	if !hasOpenQuote {
		return false
	}

	// Discover if there's any closing quote between the first quote and EOL
	hasCloseQuote, _ := l.peekAndFindAfter(`"`, openQuoteIndex)
	if !hasCloseQuote {
		return false
	}

	// Search for "\", meaning it's a multiline string
	hasBackslash, _ := l.peekAndFindAfter(`\`, openQuoteIndex)
	if !hasBackslash {
		return false
	}

	// Find next EOL that is not preceeded by an "\"
	// This means that is the last line of the multiline string
	finalEOLIndex := 0
	lastLineStartCharIndex := l.currentByteIndex + 1
	lines := 0
	for i, v := range l.input {
		if i > l.currentByteIndex {
			if v == "\n" && finalEOLIndex == 0 {
				lines += 1

				// Search backwards until the start of the line
				// to find a `\`. If none is found, this is the
				// last line of the multiline string
				lastLine := true
				for j, c := range l.input {
					if j > lastLineStartCharIndex && j < i {
						if c == `\` {
							lastLine = false
						}
					}
				}

				// Last line was found
				if lastLine {
					finalEOLIndex = i
				}

				// Update the start line of the new line
				lastLineStartCharIndex = i + 1
			}
		}
	}

	// Multiline string was not properly finished, panic
	if finalEOLIndex == 0 {
		panic("ERROR in lexer: multiline string is not closed")
	}

	// Change 4 spaces into a "\t"
	newMultilineStringValue := make([]string, len(l.input[l.currentByteIndex:finalEOLIndex]))
	copy(newMultilineStringValue, l.input[l.currentByteIndex:finalEOLIndex])
	for i, v := range newMultilineStringValue {
		if i <= len(newMultilineStringValue)-4 && v != "\n" {
			next4Characters := strings.Join(newMultilineStringValue[i:i+3], "")
			if len(strings.TrimSpace(next4Characters)) == 0 {
				arrayLeft := newMultilineStringValue[:i-1]
				arrayRight := newMultilineStringValue[i+3:]
				insertTab := append(arrayLeft, "\t")
				newMultilineStringValue = append(insertTab, arrayRight...)
			}
		}
	}

	// Create prototype and add all lines of the multiline string
	l.proto = prototype{
		kind:  keyMultiString,
		value: newMultilineStringValue,
		position: position{
			Length: finalEOLIndex - l.currentByteIndex, // This is not len(newMultilineStringValue)
			// Because the string is not being modified in the original l.input, but in a copy of it
		},
	}
	l.currentLine += lines
	l.next(false, true)
	return true
}

// Attribute type: Integer (keyInt)
func (l *lexer) lexAttrInt() bool {
	// Has to be proceeded by keyAttrDef
	if l.previousItem().kind != keyAttrDef {
		return false
	}

	// Check if this range is a possible int
	checkIntAsStr, firstIndex, lastIndex := l.searchNonWhitespacedValue()
	if checkIntAsStr == "" {
		return false
	}
	_, checkInt := strconv.Atoi(checkIntAsStr)
	if checkInt != nil {
		return false
	}

	// Create prototype and call next
	l.proto = prototype{
		kind:  keyInt,
		value: l.input[firstIndex:lastIndex],
		position: position{
			Length: lastIndex - l.currentByteIndex,
		},
	}
	l.next(false, false)
	return true
}

// Attribute type: Float (keyFloat)
func (l *lexer) lexAttrFloat() bool {
	// Has to be proceeded by keyAttrDef
	if l.previousItem().kind != keyAttrDef {
		return false
	}

	// Check if this range is a possible float
	checkFloatAsStr, firstIndex, lastIndex := l.searchNonWhitespacedValue()
	if checkFloatAsStr == "" {
		return false
	}
	_, checkInt := strconv.ParseFloat(checkFloatAsStr, 32)
	if checkInt != nil {
		return false
	}

	// Create prototype and call next
	l.proto = prototype{
		kind:  keyFloat,
		value: l.input[firstIndex:lastIndex],
		position: position{
			Length: (lastIndex - 1) - l.currentByteIndex,
		},
	}
	l.next(false, false)
	return true
}

// Attribute type: Bool (keyBool)
func (l *lexer) lexAttrBool() bool {
	// Has to be proceeded by keyAttrDef
	if l.previousItem().kind != keyAttrDef {
		return false
	}

	// Check if this range is a possible float
	checkBoolAsStr, firstIndex, lastIndex := l.searchNonWhitespacedValue()
	if checkBoolAsStr == "" {
		return false
	}
	_, checkInt := strconv.ParseBool(checkBoolAsStr)
	if checkInt != nil {
		return false
	}

	// Create prototype and call next
	l.proto = prototype{
		kind:  keyBool,
		value: l.input[firstIndex:lastIndex],
		position: position{
			Length: (lastIndex - 1) - l.currentByteIndex,
		},
	}
	l.next(false, false)
	return true
}

// Attribute type: Arithmetic operation (keyArithmetic)
func (l *lexer) lexAttrArithmetic() bool {
	// Has to be proceeded by keyAttrDef
	if l.previousItem().kind != keyAttrDef {
		return false
	}

	// Check if any arithmetic operation symbol is found
	// Possible symbols are + - * / %
	// Exponentiation (*) is not checked here because the multiplication
	// symbol will also find it
	hasArithmeticSymbol := l.peekAndFindMany([]string{"+", "-", "*", "/", "%"})
	if !hasArithmeticSymbol {
		return false
	}

	// Search for next comma, comment or EOL
	hasEndOfElem, endOfElem := l.peekAndFind(",")
	if !hasEndOfElem {
		hasEndOfElem, endOfElem = l.peekComment(l.currentByteIndex)
	}
	if !hasEndOfElem {
		endOfElem = l.peekEOL()
	}

	// Create prototype and call next
	l.proto = prototype{
		kind:  keyArithmetic,
		value: l.input[l.currentByteIndex:endOfElem],
		position: position{
			Length: (endOfElem - 1) - l.currentByteIndex,
		},
	}
	l.next(false, false)
	return true
}

// Attribute type: Comparison operation (keyComparison)
func (l *lexer) lexAttrComparison() bool {
	// Has to be proceeded by keyAttrDef
	if l.previousItem().kind != keyAttrDef {
		return false
	}

	// Check if any comparison symbol is found
	// Possible symbols are == != > >= < <=
	// By checking for > and <, it automatically checks for >= and <=
	// The same goes to checking = for == and !=
	hasComparisonSymbol := l.peekAndFindMany([]string{"=", ">", "<"})
	if !hasComparisonSymbol {
		return false
	}

	// Search for next comma, comment or EOL
	hasEndOfElem, endOfElem := l.peekAndFind(",")
	if !hasEndOfElem {
		hasEndOfElem, endOfElem = l.peekComment(l.currentByteIndex)
	}
	if !hasEndOfElem {
		endOfElem = l.peekEOL()
	}

	// Create prototype and call next
	l.proto = prototype{
		kind:  keyComparison,
		value: l.input[l.currentByteIndex:endOfElem],
		position: position{
			Length: (endOfElem - 1) - l.currentByteIndex,
		},
	}
	l.next(false, false)
	return true
}

// Attribute type: Condition operation (keyCondition)
func (l *lexer) lexAttrCondition() bool {
	// Has to be proceeded by keyAttrDef
	if l.previousItem().kind != keyAttrDef {
		return false
	}

	// Check if any condition call is found
	// Possible symbols are if for
	// Don't use peekAndFindMany here as it can only search for
	// single characters and not words
	hasConditionSymbol := false
	for i, v := range l.input {
		if i >= l.currentByteIndex && i < l.peekEOL() {
			// IF
			if v == "i" {
				if l.input[i+1] == "f" {
					hasConditionSymbol = true
					break
				}
			}
			// FOR
			if v == "f" {
				if l.input[i+1] == "o" && l.input[i+2] == "r" {
					hasConditionSymbol = true
					break
				}
			}
		}
	}
	if !hasConditionSymbol {
		return false
	}

	// Search for next comma, comment or EOL
	hasEndOfElem, endOfElem := l.peekAndFind(",")
	if !hasEndOfElem {
		hasEndOfElem, endOfElem = l.peekComment(l.currentByteIndex)
	}
	if !hasEndOfElem {
		endOfElem = l.peekEOL()
	}

	// Create prototype and call next
	l.proto = prototype{
		kind:  keyCondition,
		value: l.input[l.currentByteIndex:endOfElem],
		position: position{
			Length: (endOfElem - 1) - l.currentByteIndex,
		},
	}
	l.next(false, false)
	return true
}

// Attribute type: Function call (keyFunction)
func (l *lexer) lexAttrFunction() bool {
	// Has to be proceeded by keyAttrDef
	if l.previousItem().kind != keyAttrDef {
		return false
	}

	// Search for next comment or EOL
	hasEndOfElem, endOfElem := l.peekComment(l.currentByteIndex)
	if !hasEndOfElem {
		endOfElem = l.peekEOL()
	}

	// Search if there's any call to a function in this range
	searchFunctionCall := strings.Join(l.input[l.currentByteIndex:endOfElem], "")
	stringFunctionNames := []string{"upper(", "lower(", "append(", "concat(", "contains(", "length("}
	numericalFunctionNames := []string{"power(", "floor(", "remainder("}
	gateLogicFunctionNames := []string{"and(", "or(", "nand(", "nor(", "xor(", "xnor("}
	allFunctionNames := append(stringFunctionNames, numericalFunctionNames...)
	allFunctionNames = append(allFunctionNames, gateLogicFunctionNames...)
	if !hasPrefixToMany(searchFunctionCall, allFunctionNames) {
		return false
	}
	// Create prototype and call next
	l.proto = prototype{
		kind:  keyFunction,
		value: l.input[l.currentByteIndex:endOfElem],
		position: position{
			Length: (endOfElem - 1) - l.currentByteIndex,
		},
	}
	l.next(false, false)
	return true
}

// BLOCKS
// Start of block
func (l *lexer) lexBlockStart() bool {
	// Assume current index is an opening brace ({)
	// If it's not, peekAndFind will look and return if
	// any opening brace is found
	hasOpenBrace := true
	openBraceIndex := l.currentByteIndex
	if l.currentByte != `{` {
		hasOpenBrace, openBraceIndex = l.peekAndFind(`{`)
	}
	if !hasOpenBrace {
		return false
	}

	// Create prototype and call next
	l.proto = prototype{
		kind:  keyBlockStart,
		value: l.input[l.currentByteIndex:openBraceIndex],
		position: position{
			Length: openBraceIndex - l.currentByteIndex,
		},
	}
	l.next(false, false)
	return true
}

// End of block
func (l *lexer) lexBlockEnd() bool {
	// Assume current index is a closing brace (})
	// If it's not, peekAndFind will look and return if
	// any closing brace is found
	hasCloseBrace := true
	closeBraceIndex := l.currentByteIndex
	if l.currentByte != `}` {
		hasCloseBrace, closeBraceIndex = l.peekAndFind(`}`)
	}
	if !hasCloseBrace {
		return false
	}

	// Create prototype and call next
	l.proto = prototype{
		kind:  keyBlockEnd,
		value: l.input[l.currentByteIndex:closeBraceIndex],
		position: position{
			Length: closeBraceIndex - l.currentByteIndex,
		},
	}
	l.next(false, false)
	return true
}

// ERROR
// This means that none of the keys above were found
func (l *lexer) lexError() bool {
	// Create prototype and call next
	l.proto = prototype{
		kind:  keyError,
		value: l.input[l.currentByteIndex : l.currentByteIndex+1],
		position: position{
			Length: 1,
		},
	}
	l.next(false, false)
	return true
}

// CALL LEXERS
// Parses a CAFE file through the lexer
func (l *lexer) lexInput(debug bool) {
	for !l.atEOF {
		if debug {
			fmt.Println("DEBUG lexInput: byte: ", l.currentByte)
		}
		l.lexByte(debug)
	}
}

// Calls all lexers in a specific order to decode the
// next couple of bytes
func (l *lexer) lexByte(debug bool) {
	if debug {
		fmt.Println("DEBUG lexByte: lexEOF")
	}
	if l.lexEOF() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexTab")
	}
	if l.lexTab() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexEOL")
	}
	if l.lexEOL() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexWhitespace")
	}
	if l.lexWhitespace() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexComment")
	}
	if l.lexComment() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexAttributeDef")
	}
	if l.lexAttributeDef() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexBlockStart")
	}
	if l.lexBlockStart() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexBlockEnd")
	}
	if l.lexBlockEnd() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexAttrFunction")
	}
	if l.lexAttrFunction() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexArrayStart")
	}
	if l.lexArrayStart() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexArrayEnd")
	}
	if l.lexArrayEnd() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexArrayElem")
	}
	if l.lexArrayElem() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexAttrMultiString")
	}
	if l.lexAttrMultiString() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexAttrString")
	}
	if l.lexAttrString() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexAttrCondition")
	}
	if l.lexAttrCondition() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexAttrArithmetic")
	}
	if l.lexAttrArithmetic() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexAttrComparison")
	}
	if l.lexAttrComparison() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexAttrInt")
	}
	if l.lexAttrInt() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexAttrFloat")
	}
	if l.lexAttrFloat() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexAttrBool")
	}
	if l.lexAttrBool() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexAttrCall")
	}
	if l.lexAttrCall() {
		return
	}

	if debug {
		fmt.Println("DEBUG lexByte: lexError")
	}
	if l.lexError() {
		return
	}
}

// Transforms the int keyKind into a string
func keyKindStr(kind keyKind) string {
	switch kind {
	case keyNIL:
		return "NULL"
	case keyEOF:
		return "EOF"
	case keyTab:
		return "TAB"
	case keyError:
		return "ERROR"
	case keyComment:
		return "comment"
	case keyAttrDef:
		return "attribute definition"
	case keyBlockStart:
		return "block start"
	case keyBlockEnd:
		return "block end"
	case keyAttrCall:
		return "attribute call"
	case keyString:
		return "string"
	case keyMultiString:
		return "multiline string"
	case keyInt:
		return "int"
	case keyFloat:
		return "float"
	case keyBool:
		return "boolean"
	case keyArrayStart:
		return "array start"
	case keyArrayEnd:
		return "array end"
	case keyArrayElem:
		return "array element"
	case keyArithmetic:
		return "arithmetic operation"
	case keyComparison:
		return "comparison"
	case keyCondition:
		return "condition"
	default:
		return "unknown"
	}
}
