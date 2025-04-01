package handlers

import (
	"fmt"
	"testing"
)

type BlogQueryParamsTest struct {
	URLValues map[string]string
	Have      int
	Want      int
}

func TestParseBlogQueryParams(t *testing.T) {
	tests := []BlogQueryParamsTest{
		{
			URLValues: map[string]string{
				"test": "test",
			},
			Have: 0,
			Want: 0,
		},
	}

	fmt.Println(tests)
}
