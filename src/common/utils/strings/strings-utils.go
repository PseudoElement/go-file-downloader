package string_utils

import (
	"strings"
)

func IsEqual(str1 string, str2 string) bool {
	return strings.ToLower(str1) == strings.ToLower(str2)
}
