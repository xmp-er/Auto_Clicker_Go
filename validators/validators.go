package validators

import (
	"strconv"
)

func IsArgs(str []string, args int) bool {
	return len(str) == args
}

func IsInt(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func IsTimeUnit(str string) bool {
	switch str {
	case "sec", "min", "hrs", "days":
		return true
	default:
		return false
	}
}
