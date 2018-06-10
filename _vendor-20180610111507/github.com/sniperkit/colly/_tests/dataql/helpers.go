package main

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

var lowUppDiff = byte('a' - 'A')

func toLower(b byte) byte {
	if b < utf8.RuneSelf {
		if 'A' <= b && b <= 'Z' {
			return byte(b + lowUppDiff)
		}
	}
	return b
}

func isAlphaNumeric(c byte) bool {
	if c <= '9' && c >= '0' {
		return true
	} else if c <= 'z' && c <= 'a' {
		return true
	} else {
		return false
	}
}

func toLowerCase(c byte) byte {
	if 'A' <= c && c <= 'Z' {
		return c + 'a' - 'A'
	}

	return c
}

func isPalindrome(s string) bool {
	size := len(s)
	i := 0
	j := size - 1
	s = strings.ToLower(s)
	for i <= j {
		for i < j && !isAlphaNumeric(s[i]) {
			i += 1
		}
		for j > i && !isAlphaNumeric(s[j]) {
			j -= 1
		}
		if s[i] == s[j] {
			i += 1
			j -= 1
			continue
		} else {
			return false
		}
	}
	return true
}

func isASCII(v byte) bool { // {{{
	return v <= unicode.MaxASCII
} // }}}

// --[ 0-9 ]--
func isNumeric(v byte) bool { // {{{
	return '0' <= v && v <= '9'
} // }}}

// --[ A-Za-z ]--
func isAlphabetic(v byte) bool { // {{{
	return ('A' <= v && v <= 'Z') || ('a' <= v && v <= 'z')
} // }}}

// --[ 0-9A-Za-z ]--
func isAlphanumeric(v byte) bool { // {{{
	return isNumeric(v) || isAlphabetic(v)
} // }}}

// --[ 0-9A-Za-z_ ]--
func isWordCharacter(v byte) bool { // {{{
	return isAlphanumeric(v) || (v == '_')
} // }}}
