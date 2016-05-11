package actions

import (
	"strings"
)

func commitFocus(str string) string {
	// first, some sanitization
	// parse only the title, strip the commit body
	str = strings.Split(str, "\n")[0]
	if !strings.Contains(str, "(") || !strings.Contains(str, ")") {
		return ""
	}
	// fetch the string between the parentheses
	return strings.TrimSpace(strings.Split(strings.Split(str, ")")[0], "(")[1])
}

func commitTitle(str string) string {
	// first, some sanitization
	// parse only the title, strip the commit body
	str = strings.Split(str, "\n")[0]
	// if the commit title doesn't follow our standards, just dump the whole string
	if !strings.Contains(str, ":") {
		return str
	}
	return strings.TrimSpace(strings.Split(str, ":")[1])
}
