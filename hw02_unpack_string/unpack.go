package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputStr string) (string, error) {
	if inputStr == "" {
		return "", nil
	}

	var sb strings.Builder

	toSkip := false

	var prevRune rune

	for idx, r := range inputStr {
		if idx == 0 {
			if unicode.IsDigit(r) {
				return "", ErrInvalidString
			}
		}

		if r == '\\' {
			toSkip = true
			continue
		}

		if toSkip {
			sb.WriteRune(r)
			prevRune = 0
			toSkip = false
			continue
		}

		if unicode.IsDigit(r) {
			if prevRune == 0 {
				return "", ErrInvalidString
			}
			count, err := strconv.Atoi(string(r))
			if err != nil {
				return "", ErrInvalidString
			}
			if count == 0 {
				DeleteLastRune(&sb)
			}
			for ; count > 1; count-- {
				sb.WriteRune(prevRune)
			}
			prevRune = 0
		} else {
			sb.WriteRune(r)
			prevRune = r
		}
	}
	return sb.String(), nil
}

func DeleteLastRune(sb *strings.Builder) {
	str := sb.String()
	if len(str) > 0 {
		str = str[:len(str)-1]
	}
	sb.Reset()
	sb.WriteString(str)
}
