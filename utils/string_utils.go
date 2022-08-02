package utils

import (
	"strings"
)

func Head(text []rune, length int) string {
	if len(text) == 0 || len(text) <= length {
		return string(text)
	}
	if length <= 0 {
		return ""
	}
	return string(text[:length])
}

func Tail(text []rune, length int) []rune {
	if len(text) == 0 {
		return text
	}
	if length <= 0 {
		return []rune{}
	}
	return text[len(text)-length:]
}

func Take(text []rune, begin int) []rune {
	if len(text) == 0 {
		return text
	}
	if begin <= 0 {
		begin = 0
	}
	if begin > len(text)-1 {
		return []rune{}
	}

	return text[begin:]
}

func Substring(text []rune, begin, end int) []rune {
	if len(text) == 0 {
		return text
	}
	if begin < 0 {
		begin = 0
	}
	if end >= len(text)-1 {
		end = len(text) - 1
	}
	if begin > end {
		return []rune{}
	}

	return text[begin : end+1]
}

func IsAnsiChars(text string) bool {
	if len(text) == 0 {
		return false
	}
	for _, v := range []rune(text) {
		if !((v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z')) {
			return false
		}
	}
	return true
}

func IsNumericChars(text string) bool {
	if len(text) == 0 {
		return false
	}
	for _, v := range []rune(text) {
		if v < '0' || v > '9' {
			return false
		}
	}
	return true
}

func Remove(text []rune, chars []rune, exclude string) string {
	if len(text) == 0 || len(chars) == 0 {
		return string(text)
	}
	charSet := make(map[rune]struct{})
	for _, v := range chars {
		charSet[v] = struct{}{}
	}
	var sb strings.Builder
	var removed bool
	for _, v := range text {
		_, ok := charSet[v]
		if ok && !strings.Contains(exclude, string(v)) {
			removed = true
			continue
		}
		sb.WriteRune(v)
	}
	if removed {
		return sb.String()
	} else {
		return string(text)
	}
}

func RemoveRepeatNum(text []rune, n int) string {
	if len(text) < n {
		return string(text)
	}
	var sb strings.Builder
	var cnt int
	for i, v := range text {
		if v >= '0' && v <= '9' {
			cnt++
			continue
		}
		if cnt > 0 && cnt < n {
			sb.WriteString(string(text[i-cnt : i]))
		}
		cnt = 0
		sb.WriteRune(v)
	}
	if cnt > 0 && cnt < n {
		sb.WriteString(string(text[len(text)-cnt:]))
	}

	return sb.String()
}
