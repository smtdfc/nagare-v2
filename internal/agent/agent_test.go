package agent

import (
	"testing"

	"github.com/smtdfc/nagare/pkgs/logger"
	"github.com/stretchr/testify/assert"
)

func TestAgent_WithState(t *testing.T) {
	agent := NewAgent(logger.NewLogger())
	state := NewState()

	assert.Nil(t, agent.State)
	agent.WithState(state)
	assert.NotNil(t, agent.State)
}

func TestNewAgent(t *testing.T) {
	agent := NewAgent(logger.NewLogger())
	assert.NotNil(t, agent)
	assert.Nil(t, agent.State)
	assert.NotNil(t, agent.logger)
}

func TestNewAgent_DefaultLogger(t *testing.T) {
	agent := NewAgent(nil)
	assert.NotNil(t, agent.logger)
}
