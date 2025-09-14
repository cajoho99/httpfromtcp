package headers

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	bef, _, found := bytes.Cut(data, []byte("\r\n"))

	if !found {
		return 0, false, nil
	}

	if len(bef) == 0 {
		return 2, true, nil
	}

	key, value, found := bytes.Cut(bef, []byte(":"))

	if !found {
		return 0, false, errors.New(": not found when parsing")
	}

	if len(value) < 1 {
		return 0, false, errors.New("Value is empty")
	}

	strValue := strings.TrimSpace(string(value))

	strKey, err := ValidKey(key)

	if err != nil {
		return 0, false, err
	}

	h[strKey] = strValue

	// + 2 for the CLRF
	return len(bef) + 2, false, nil
}

func (h Headers) Get(key string) (val string) {
	lKey := strings.ToLower(key)
	return h[lKey]
}

func ValidKey(data []byte) (s string, err error) {
	if len(data) < 1 {
		return "", errors.New("Value is empty")
	}

	if unicode.IsSpace(rune(data[len(data)-1])) {
		return "", errors.New("Whitespace before : is not allowed")
	}

	strKey := string(data)
	strKey = strings.TrimSpace(strKey)
	strKey = strings.ToLower(strKey)

	if !hasValidChars(strKey) {
		return "", errors.New("Has invalid chars")
	}

	return strKey, nil
}

func hasValidChars(s string) bool {
	validSpecialChars := []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}
	for _, r := range s {
		if !(unicode.IsLetter(r) || unicode.IsNumber(r) || slices.Contains(validSpecialChars, r)) {
			return false
		}
	}
	return true
}
