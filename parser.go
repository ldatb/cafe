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
	"strings"
)

// attrKind  defines all kinds of possible Attributes of an CAFE file
type attrKind int

const (
	attrNIL        attrKind = iota // 0
	attrString                     // 1
	attrInt                        // 2
	attrFloat                      // 3
	attrBool                       // 4
	attrArray                      // 5
	attrArithmetic                 // 6
	attrComparison                 // 7
	attrCondition                  // 8
	attrFunction                   // 9
)

// The Parser is called after the input goes through the lexer
// It will fetch all the elements and put them into the respective
// place
type Parser struct {
	// The lexer
	lx *lexer

	// Current item in the items list
	currentItem item

	// Index of the current item in the items list
	currentItemIndex int

	// List of global Attributes
	Attributes map[string]attribute

	// List of the current nested Blocks names
	currentBlocks []string

	// List of Blocks of the file
	Blocks map[string]block

	// Last item was reached
	atLastItem bool
}

// attribute defines the variables of an CAFE file
type attribute struct {
	// Name of the attribute
	Name string

	// Attribute value
	Value interface{}

	// Kind of the attribute
	kind attrKind
}

// Blocks are structures in an CAFE that can hold multiple
// Attributes and even other Blocks
type block struct {
	// Name of the block
	Name string

	// List of Attributes of the block
	Attributes map[string]attribute

	// List of Blocks inside this block
	Blocks map[string]block
}

// Creates a Parser
func newParser(input []string) *Parser {
	lx := newLexer(input)
	lx.lexInput(false)

	return &Parser{
		lx:               lx,
		currentItem:      lx.items[0],
		currentItemIndex: 0,
		Attributes:       map[string]attribute{},
		currentBlocks:    []string{},
		Blocks:           map[string]block{},
		atLastItem:       false,
	}
}

// Moves the Parser to the next item
func (p *Parser) nextItem(nextCount int) {
	p.currentItemIndex += nextCount

	if p.currentItemIndex == len(p.lx.items) {
		p.atLastItem = true
		return
	}

	p.currentItem = p.lx.items[p.currentItemIndex]
}

// Returns the next item of the current one
// Does not move the Parser to the next item
func (p *Parser) peekNextItem() item {
	if p.currentItemIndex == len(p.lx.items)-1 {
		return item{}
	}
	return p.lx.items[p.currentItemIndex+1]
}

// Goes recurssively into the nested Blocks until
// it reaches the last one in the nest
// Returns a boolean (is in block = true / is not in block = false)
// Returns a pointer to the block (if any)
func (p *Parser) getCurrentBlock() (bool, *block) {
	if len(p.currentBlocks) == 0 {
		return false, nil
	}

	// Get first block
	lastBlock := p.Blocks[p.currentBlocks[0]]

	// Go recursively into the Blocks and find the last one
	for i, b := range p.currentBlocks {
		if i == 0 {
			continue
		}
		lastBlock = lastBlock.Blocks[b]
	}

	return true, &lastBlock
}

// Transforms a keyKind in an attrKind
func keyKindToAttrKind(k keyKind) attrKind {
	switch k {
	case keyString:
		return attrString
	case keyMultiString:
		return attrString
	case keyInt:
		return attrInt
	case keyFloat:
		return attrFloat
	case keyBool:
		return attrBool
	case keyArrayStart:
		return attrArray
	case keyArrayElem:
		return attrArray
	case keyArithmetic:
		return attrArithmetic
	case keyComparison:
		return attrComparison
	case keyCondition:
		return attrCondition
	case keyFunction:
		return attrFunction
	default:
		return attrNIL
	}
}

// Parses an item into an attribute
func (p *Parser) parseAttribute() bool {
	if p.currentItem.kind != keyAttrDef {
		return false
	}

	// Get attribute value and kind
	itemItem := p.peekNextItem()
	itemvalue := itemItem.value
	nextCount := 2

	// Item is an array
	if itemItem.kind == keyArrayStart {
		itemvalue = ""
		// Get items until keyArrayEnd
		arrayItems := []string{}
		for i, v := range p.lx.items {
			if i > p.currentItemIndex {
				if v.kind == keyArrayEnd {
					break
				}
				arrayItems = append(arrayItems, v.value)
				nextCount += 1
			}
		}
		itemvalue = strings.Join(arrayItems, ", ")
	}

	// Transform value string into interface
	attrvalue := transformItem(itemvalue, itemItem.kind)

	// Build attribute
	newAttr := attribute{
		Name:  p.currentItem.value,
		Value: attrvalue,
		kind:  keyKindToAttrKind(itemItem.kind),
	}

	// Add new attribute into global or nested block
	isBlock, currentBlock := p.getCurrentBlock()
	if isBlock {
		currentBlock.Attributes[newAttr.Name] = newAttr
	} else {
		p.Attributes[newAttr.Name] = newAttr
	}

	// Call next item and return
	p.nextItem(nextCount)
	return true
}

// Skip array elements since these are already being checked
// by parseAttribute
func (p *Parser) parseArrayElement() bool {
	if p.currentItem.kind != keyArrayElem {
		return false
	}
	p.nextItem(1)
	return true
}

// Parses a block start
func (p *Parser) parseBlockStart() bool {
	if p.currentItem.kind != keyBlockStart {
		return false
	}

	// Build block
	newBlock := block{
		Name:       p.currentItem.value,
		Attributes: map[string]attribute{},
		Blocks:     map[string]block{},
	}

	// Add new block into global or nested block
	isBlock, currentBlock := p.getCurrentBlock()
	if isBlock {
		currentBlock.Blocks[newBlock.Name] = newBlock
	} else {
		p.Blocks[newBlock.Name] = newBlock
	}

	// Add block name to the array of current Blocks
	p.currentBlocks = append(p.currentBlocks, newBlock.Name)

	// Call next item and return
	p.nextItem(1)
	return true
}

// Parses a block end
func (p *Parser) parseBlockEnd() bool {
	if p.currentItem.kind != keyBlockEnd {
		return false
	}

	// Remove last block name of the array of current Blocks
	p.currentBlocks = p.currentBlocks[:len(p.currentBlocks)-1]

	// Call next and return
	p.nextItem(1)
	return true
}

// Parse EOF
func (p *Parser) parseEOF() bool {
	if p.currentItem.kind != keyEOF {
		return false
	}
	p.nextItem(1)
	return true
}

// Parse others: comment, EOL, NIL, ERROR
func (p *Parser) parseOthers() bool {
	if p.currentItem.kind != keyComment && p.currentItem.kind != keyNIL && p.currentItem.kind != keyError {
		return false
	}
	p.nextItem(1)
	return true
}

// Parse all items until the last one
func (p *Parser) parseItems(debug bool) {
	for !p.atLastItem {
		if debug {
			fmt.Println("DEBUG ITEM:", p.currentItem.value)
		}
		p.parseItem(debug)
	}
}

// Calls all Parsers in a specific order to parse the next item
func (p *Parser) parseItem(debug bool) {
	if debug {
		fmt.Println("DEBUG parseItem: parseEOF")
	}
	if p.parseEOF() {
		return
	}

	if debug {
		fmt.Println("DEBUG parseItem: parseAttribute")
	}
	if p.parseAttribute() {
		return
	}

	if debug {
		fmt.Println("DEBUG parseItem: parseArrayElement")
	}
	if p.parseArrayElement() {
		return
	}

	if debug {
		fmt.Println("DEBUG parseItem: parseBlockStart")
	}
	if p.parseBlockStart() {
		return
	}

	if debug {
		fmt.Println("DEBUG parseItem: parseBlockEnd")
	}
	if p.parseBlockEnd() {
		return
	}

	if debug {
		fmt.Println("DEBUG parseItem: parseOthers")
	}
	if p.parseOthers() {
		return
	}

	// No value found so skip to next item
	if debug {
		fmt.Println("DEBUG parseItem: Skipping to next item")
	}
	p.nextItem(1)
}
