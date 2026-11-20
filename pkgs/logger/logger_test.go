package logger_test

import (
	"testing"

	"github.com/smtdfc/nagare/pkgs/logger"
	"github.com/stretchr/testify/assert"
)

func Test_NewLogger(t *testing.T) {
	l := logger.NewLogger()

	assert.NotNil(t, l.Logger, "logger not nil")
}
