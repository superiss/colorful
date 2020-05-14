package colorful

import (
	"bytes"
	"fmt"
	"log"
)

const (
	escapeChar = '@'       // Escape character for color syntax
	resetCode  = "\033[0m" // Short for reset to default style
)

// Mapping from character to concrete escape code.
var codeMap = map[int]int{
	'|': 0,
	'!': 1,
	'.': 2,
	'/': 3,
	'_': 4,
	'^': 5,
	'&': 6,
	'?': 7,
	'-': 8,
	'*': 60,

	'k': 30,
	'r': 31,
	'g': 32,
	'y': 33,
	'b': 34,
	'm': 35,
	'c': 36,
	'w': 37,
	'd': 39,

	'K': 40,
	'R': 41,
	'G': 42,
	'Y': 43,
	'B': 44,
	'M': 45,
	'C': 46,
	'W': 47,
	'D': 49,
}

// Colorize Compile color syntax string like "rG" to escape code.
func colorize(x string) string {
	attr := 0
	fg := 39
	bg := 49

	for _, key := range x {
		c, ok := codeMap[int(key)]
		switch {
		case !ok:
			log.Printf("Wrong color syntax: %c", key)
		case 0 <= c && c <= 8:
			attr = c
		case 30 <= c && c <= 37:
			fg = c
		case 40 <= c && c <= 47:
			bg = c
		case c == 60:
			fg += c
		}
	}
	return fmt.Sprintf("\033[%d;%d;%dm", attr, fg, bg)
}

// Handle state after meeting one '@'
func compileColorSyntax(input, output *bytes.Buffer) {
	i, _, err := input.ReadRune()
	if err != nil {
		// EOF got
		log.Print("Parse failed on color syntax")
		return
	}

	switch i {
	default:
		output.WriteString(colorize(string(i)))
	case '{':
		color := bytes.NewBufferString("")
		for {
			i, _, err := input.ReadRune()
			if err != nil {
				log.Print("Parse failed on color syntax")
				break
			}
			if i == '}' {
				break
			}
			color.WriteRune(i)
		}
		output.WriteString(colorize(color.String()))
	case escapeChar:
		output.WriteRune(escapeChar)
	}
}

// Compile the string and replace color syntax with concrete escape code.
func compile(x string) string {
	if x == "" {
		return ""
	}

	input := bytes.NewBufferString(x)
	output := bytes.NewBufferString("")

	for {
		i, _, err := input.ReadRune()
		if err != nil {
			break
		}
		switch i {
		default:
			output.WriteRune(i)
		case escapeChar:
			compileColorSyntax(input, output)
		}
	}
	return output.String()
}

// Compile multiple values, only do compiling on string type.
func compileValues(a *[]interface{}) {
	for i, x := range *a {
		if str, ok := x.(string); ok {
			(*a)[i] = compile(str)
		}
	}
}

// Sprint Similar to fmt.Print, will reset the color at the end.
func Sprint(format string, a ...interface{}) string {
	format += resetCode
	format = compile(format)
	return fmt.Sprintf(format, a...)
}
