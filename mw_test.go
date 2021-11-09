package mw_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/dcmn-com/mw"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestTracing_GenerateRequestID(t *testing.T) {
    var numCalls int

    h := mw.Tracing(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        ctx := req.Context()

        id, ok := mw.ContextRequestID(ctx)
        require.True(t, ok)

        assert.NotEmpty(t, id)
        assert.Equal(t, w.Header().Get("X-Request-Id"), id)

        numCalls++
        w.WriteHeader(http.StatusNoContent)
    }))

    req := httptest.NewRequest("GET", "http://example.org", nil)
    rec := httptest.NewRecorder()

    h.ServeHTTP(rec, req)
    require.Equal(t, 1, numCalls)
}

func TestTracing_ForwardRequestID(t *testing.T) {
    reqID := "abc"

    var numCalls int

    h := mw.Tracing(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        ctx := req.Context()

        id, ok := mw.ContextRequestID(ctx)
        require.True(t, ok)

        assert.Equal(t, reqID, id)
        assert.Equal(t, w.Header().Get("X-Request-Id"), id)

        numCalls++
        w.WriteHeader(http.StatusNoContent)
    }))

    req := httptest.NewRequest("GET", "http://example.org", nil)
    req.Header.Set("X-Request-Id", reqID)

    rec := httptest.NewRecorder()

    h.ServeHTTP(rec, req)
    require.Equal(t, 1, numCalls)
}
