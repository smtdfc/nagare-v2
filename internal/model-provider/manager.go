package modelprovider

import "github.com/smtdfc/nagare/pkgs/logger"

type Manager struct {
	logger *logger.Logger
}

// @Injectable
func NewModelProviderManager(logger *logger.Logger) *Manager {
	return &Manager{
		logger: logger,
	}
}
