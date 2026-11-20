package message_test

import (
	"testing"

	"github.com/smtdfc/nagare/pkgs/message"
	"github.com/stretchr/testify/assert"
)

func Test_NewAgentStartedMessage(t *testing.T) {
	m := message.NewAgentStartedMessage()
	assert.NotNil(t, m, "NewAgentStartedMessage() should not be nil")
	assert.Equal(t, message.AgentStarted, m.Type, "NewAgentStartedMessage() should have the correct type")
}

func Test_NewAgentEndedMessage(t *testing.T) {
	m := message.NewAgentEndedMessage(
		false,
		10,
		"error",
	)

	assert.NotNil(t, m, "NewAgentEndedMessage() should not be nil")
	assert.Equal(t, message.AgentEnded, m.Type, "NewAgentEndedMessage() should have the correct type")
}
