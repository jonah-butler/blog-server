package blog

import (
	"strings"
)

func generateSlug(title string) string {
	return strings.ToLower(strings.Join(strings.Split(title, " "), "-"))
}
