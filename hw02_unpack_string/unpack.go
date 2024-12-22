package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func IsLetterOrNewString(r rune) bool {
	return unicode.IsLetter(r) || string(r) == "\n"
}

func IsContainAllowedSymbols(inputStr string) bool {
	for _, r := range inputStr {
		isLetter := IsLetterOrNewString(r)
		isDigit := unicode.IsDigit(r)
		if !isLetter && !isDigit {
			return false
		}
	}
	return true
}

func Unpack(inputStr string) (string, error) {
	if inputStr == "" {
		return "", nil
	}
	if !IsContainAllowedSymbols(inputStr) {
		return "", ErrInvalidString
	}

	var sb strings.Builder

	var curSymb, prevSymb rune

	for idx, r := range inputStr {
		if idx == 0 {
			if unicode.IsDigit(r) {
				return "", ErrInvalidString
			}
			prevSymb = r
			continue
		}

		curSymb = r
		switch {
		case IsLetterOrNewString(prevSymb) && IsLetterOrNewString(curSymb):
			sb.WriteString(string(prevSymb))
			prevSymb = curSymb
		case IsLetterOrNewString(prevSymb) && unicode.IsDigit(curSymb):
			count, err := strconv.Atoi(string(curSymb))
			if err != nil {
				return "", ErrInvalidString
			}
			for ; count > 0; count-- {
				sb.WriteString(string(prevSymb))
			}
			prevSymb = curSymb
		case unicode.IsDigit(prevSymb) && IsLetterOrNewString(curSymb):
			prevSymb = curSymb
		case unicode.IsDigit(prevSymb) && unicode.IsDigit(curSymb):
			return "", ErrInvalidString
		}
	}

	if IsLetterOrNewString(curSymb) {
		sb.WriteString(string(curSymb))
	}

	return sb.String(), nil
}
