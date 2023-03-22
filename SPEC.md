# CAFE Syntax Specification

This is the specification of the syntax and semantics of CAFE.

## Syntax Notation

This notation is intended for human consumption rather than machine consumption, with the following conventions:

- Double and single quotes (`"` and `'`) are used to mark literal character sequences, which may be either punctuation markers or keywords.
- The symbol `|` indicates that any one of its left and right operands may be present.
- Parentheses `(` and `)` are used to group items together to apply the `|` operator to them collectively.
- The symbol `=` is used to declare variables

## Lexical Elements

### Comments

Comments start with the // sequence and end with the next newline sequence. A line comment is considered equivalent to a newline sequence.

Inline comments are also supported.

### Operators and Delimiters

The following character sequences represent operators, delimiters, and other special tokens:

```
+   {   ==  <   &&  
-   }   !=  >   ||  
*   [   =   <=  !
/   ]   :   >=  (
%   ${  ?   \   )
```

## Structural elements

The structural language consists of syntax representing the following constructs:

- Attributes, which assign a value to a specified name.
- Blocks, which create a child body.
- Body Content, which consists of a collection of attributes and blocks.

```
attribute = "value"
block {
    sub_block {
        body_content = "foo"
    }
}
```

Note: Blocks **MUST** have a name assigned to it

## Data Types

CAFE supports the common data types:

- Number (assigned integers and floats): `number = 3.14` or `number = -10`
- String (a collection of characters): `string = "Hello World!`
- Multiline string (a collection of lines): 
```
multilineString = "line1" \
                    "line2" \
                    "lineN" \
                    "final line"
```
- Boolean (true of false values): `bool = false` or `bool = true`
- Array (collection of data) = `array = ["foo", "bar", 2023, false]`
- Multiline array:
```
multilineArray = [
    "this",
    "is",
    "a",
    "multiline",
    "array"
]
```
- Time (ISO 8601): `today = 2023-03-14`

Strings (and multiline strings) also support interpolation:
```
var1 = "Hello"
var2 = "${var1}, world!"
```

## Expressions

### If

A "if" is a conditional construct to make an attribute based on a condition, applying it's value by using the `?` and `:` operators.

If conditions can return all kinds of attribute types.

```
msg = "Hello World!"
has_hello = if contains(msg, "hello") ? true : false // True
has_hi = if contains(msg, "hi") ? true : false // False

// Note:
// This is only an example, in this case it would make more sense to simply do
// has_hello = contains(msg, "hello")
// has_hi = contains(msg, "hi")
```

### For

A "for loop" is a construct for constructing a collection by projecting the items from another collection. 2 variables are automatically declared in a for loop: the index and the value, you can use both to call functions, expressions, etc.

```
months = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]
monthNumber = for months : index + 1
// monthNumber = [1, 2, 3, 4, 5, 6, 7, 8, 10, 11, 12]
```

### Operations

Operations apply a particular operator to either one or more expression terms.

#### Arithmetic operators
```
a + b   // addition
a - b   // subtraction
a * b   // multiplication
a / b   // division
```

#### Comparative operators

Comparative operators are only valid for numerical attributes. To compare booleans, use the gate logic functions.

```
a == b    // Equal
a != b    // Not equal
a < b     // less than
a <= b    // less than or equal to
a > b     // greater than
a >= b    // greater than or equal to
```

### Functions

The following functions come by default with the CAFE interpreter:

#### Strings

String functions are normal functions that can modify strings

- upper(str) // Uppercases a string
- lower(str) // Lowercases a string
- append(str, val) // Adds a string to the end of another string
- concat(arr, separator) // Concatenates an array into a string. All elements become strings
- contains(str, substr) // Checks if a string contains a substring
- length(str) // Checks the length of the string

#### Numerical

Numerical functions perform calculations with floating points, but, if any of your values is an integer,
it will return the result as an integer too.

- power(value, exponent) // Exponential
- floor(dividend, divisor) // Floor division
- remainder(dividend, divisor) // Remainder of division

#### Gate Logic

Gate logic functions can only be applied to boolean values.

- and(cond1, cond2) // AND gate
- or(cond1, cond2) // OR gate
- nand(cond1, cond2) // NAND gate
- nor(cond1, cond2) // NOR gate
- xor(cond1, cond2) // XOR gate
- xnor(cond1, cond2) // XNOR gate
