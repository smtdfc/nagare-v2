package logger

import "github.com/smtdfc/nagare/pkgs/logger"

// @Injectable
func LoggerProvider() *logger.Logger {
	return logger.NewLogger()
}
