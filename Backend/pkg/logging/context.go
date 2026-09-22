package logging

import (
	"context"

	"github.com/MohammadrezaNadirkhanloo/todo-list-FullStack/pkg/appctx"
)


func fieldsFromContext(ctx context.Context) []Field {
	if ctx == nil {
		return nil
	}

	var fields []Field
	if id, ok := appctx.RequestID(ctx); ok {
		fields = append(fields, Field{Key: "request_id", Value: id})
	}
	if id, ok := appctx.UserID(ctx); ok {
		fields = append(fields, Field{Key: "user_id", Value: id})
	}
	return fields
}
