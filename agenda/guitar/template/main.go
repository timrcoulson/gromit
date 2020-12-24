package template

import (
	"strconv"
	"unicode"
)

func Print(root int, temp string) string{
	output := ""

	killNextChar := false
	for _, r := range []rune(temp) {
		if killNextChar {
			killNextChar = false
			continue
		}
		if unicode.IsNumber(r) {
			fretNumber, _ := strconv.Atoi(string(r))

			// If double char, kill the next one.
			if fretNumber + root > 9 {
				killNextChar = true
			}

			output += strconv.Itoa(fretNumber + root)
		} else {
			output += string(r)
		}
	}

	return output
}