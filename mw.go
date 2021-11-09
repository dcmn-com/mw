package mw

import (
    "net/http"
    "os"
    "time"

    "github.com/dcmn-com/jlo"
    uuid "github.com/satori/go.uuid"
)

const (
    // fieldKeyRequestID is the request ID log field name.
    fieldKeyRequestID = "@request_id"
    // RequestIDHeader defines the header, which is used to trace frontend requests
    // through multiple services.
    RequestIDHeader = "X-Request-Id"
)

// Logger adds a logger to request context and initializes it with the passed in
// commit hash and the request ID from context.
func Logger(commit string) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := r.Context()

            l := jlo.NewLogger(os.Stdout)
            l = l.WithField(jlo.FieldKeyCommit, commit)

            if requestID, _ := ContextRequestID(ctx); requestID != "" {
                l = l.WithField(fieldKeyRequestID, requestID)
            }

            ctx = ContextLoggerSet(ctx, l)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// RequestLogging logs incoming request along with the time it took to handle it.
func RequestLogging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        l := ContextLogger(ctx)

        start := time.Now()
        wrapped := NewResponseWriter(w)
        next.ServeHTTP(wrapped, r)

        l.Debugf(
            "%s %s %d %s",
            r.Method,
            r.RequestURI,
            wrapped.Status(),
            time.Since(start),
        )
    })
}

// Tracing sets request ID on the request context.
func Tracing(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        requestID := r.Header.Get(RequestIDHeader)
        if requestID == "" {
            requestID = uuid.NewV4().String()
        }

        ctx = ContextRequestIDSet(ctx, requestID)
        w.Header().Set(RequestIDHeader, requestID)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// ResponseWriter wraps a standard http.ResponseWriter to store the status code for logging.
type ResponseWriter struct {
    length int64
    status int
    http.ResponseWriter
}

// NewResponseWriter returns a wrapped response writer.
func NewResponseWriter(res http.ResponseWriter) *ResponseWriter {
    // Defaults the status code to 200
    return &ResponseWriter{0, 200, res}
}

// Header wraps http.ResponseWriter Header() method.
func (w *ResponseWriter) Header() http.Header {
    return w.ResponseWriter.Header()
}

// Length wraps http.ResponseWriter Length() method.
func (w *ResponseWriter) Length() int64 {
    return w.length
}

// Status wraps http.ResponseWriter Status() method.
func (w *ResponseWriter) Status() int {
    return w.status
}

// Write wraps http.ResponseWriter Write() method.
func (w *ResponseWriter) Write(data []byte) (int, error) {
    n, err := w.ResponseWriter.Write(data)

    w.length += int64(n)

    return n, err
}

// WriteHeader wraps http.ResponseWriter WriteHeader() method.
func (w *ResponseWriter) WriteHeader(statusCode int) {
    // Store the status code
    w.status = statusCode

    // Write the status code onward
    w.ResponseWriter.WriteHeader(statusCode)
}
