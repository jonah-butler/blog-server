package serviceutilities

import (
	ck "blog-api/contextkeys"
	"context"
)

func GetAuthorID(ctx context.Context) (string, bool) {
	authorID, ok := ctx.Value(ck.UserIDKey).(string)
	return authorID, ok
}
