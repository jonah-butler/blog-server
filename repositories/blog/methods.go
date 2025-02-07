package blog

import (
	"strings"
)

func splitAndTrim(input string) []string {
	inputSlice := strings.Split(input, ",")

	for i, category := range inputSlice {
		inputSlice[i] = strings.TrimSpace(category)
	}

	return inputSlice
}
