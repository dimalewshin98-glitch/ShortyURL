package contextkeys

import (
	"context"
	"errors"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

func GetUserID(ctx context.Context) (int, error) {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return 0, errors.New("userID not found in context")
	}
	userID, ok := val.(int)
	if !ok {
		return 0, errors.New("invalid userID type in context")
	}
	return userID, nil
}
