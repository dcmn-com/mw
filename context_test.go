package mw_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/dcmn-com/mw"

	"github.com/dcmn-com/jlo"
	"github.com/stretchr/testify/assert"
)

func TestContextSetLogger(t *testing.T) {
	testLogger := jlo.NewLogger(ioutil.Discard).WithField("request_id", "12345678-1234-1234-1234-123456789012")

	// set logger
	ctx := mw.ContextLoggerSet(context.Background(), testLogger)

	// retrieve logger
	l := mw.ContextLogger(ctx)

	assert.Equal(t, testLogger, l)
}

func TestContextSetRequestID(t *testing.T) {
	testID := "12345678-1234-1234-1234-123456789012"
	ctx := context.Background()

	// set id
	ctx = mw.ContextRequestIDSet(ctx, testID)

	// retrieve id
	id, ok := mw.ContextRequestID(ctx)

	assert.True(t, ok)
	assert.Equal(t, testID, id)
}
