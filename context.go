package mw

import (
	"context"

	"github.com/dcmn-com/jlo"
)

type contextKey int

const (
	loggerKey contextKey = iota
	requestIDKey
)

// ContextLoggerSet sets the logger on a context, used to have a request scoped logger
// which includes request ID as a field.
func ContextLoggerSet(ctx context.Context, l *jlo.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// ContextLogger retrieved the logger from a context, used in the handlers to retrieve
// a request scoped logger which includes request ID as a field from the request context.
func ContextLogger(ctx context.Context) *jlo.Logger {
	l, ok := ctx.Value(loggerKey).(*jlo.Logger)
	if !ok {
		return jlo.DefaultLogger()
	}

	return l
}

// ContextRequestIDSet sets the request ID on a context to include a unique request ID
// into logs and error responses.
func ContextRequestIDSet(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// ContextRequestID retrieve the request ID from a context. Used in the handlers to
// get the request ID e.g. to include it into error responses.
func ContextRequestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}
