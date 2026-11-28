package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggerProvider(t *testing.T) {
	log := LoggerProvider()

	require.NotNil(t, log)
	assert.NotNil(t, log.Logger)
}
