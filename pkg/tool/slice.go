package tool

import "strings"

func ContainStr(s []string, sub string) bool {
	for _, v := range s {
		if v == sub {
			return true
		}
	}
	return false
}

func PrefixStr(s []string, sub string) bool {
	for _, v := range s {
		if strings.HasPrefix(v, sub) {
			return true
		}
	}
	return false
}
