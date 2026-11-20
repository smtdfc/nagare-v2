package agent

import (
	"context"
	"testing"

	error_utils "github.com/smtdfc/nagare/core/error-utils"
	"github.com/smtdfc/nagare/core/execution"
	"github.com/smtdfc/nagare/pkgs/logger"
	"github.com/smtdfc/nagare/pkgs/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgent_Invoke_RequiresState(t *testing.T) {
	agent := NewAgent(logger.NewLogger())

	channel, err := agent.Invoke(execution.NewContext(context.Background()), &fakeAdapter{}, nil, &InvokeConfig{})

	assert.Nil(t, channel)
	assert.Equal(t, error_utils.ErrAgentStateEmpty, err)
}

func TestAgent_Invoke_RequiresAdapter(t *testing.T) {
	agent := NewAgent(logger.NewLogger()).WithState(NewState())

	channel, err := agent.Invoke(execution.NewContext(context.Background()), nil, nil, &InvokeConfig{})

	assert.Nil(t, channel)
	assert.Equal(t, error_utils.ErrAdapterNotFound, err)
}

func TestAgent_Invoke_StreamsUntilDone(t *testing.T) {
	state := NewState()
	adapter := &fakeAdapter{
		turns: [][]message.Message{
			{message.NewTextMessage("answer", message.Agent)},
		},
	}
	agent := NewAgent(logger.NewLogger()).WithState(state)
	inputs := []message.Message{message.NewTextMessage("question", message.User)}

	channel, err := agent.Invoke(
		execution.NewContext(context.Background()),
		adapter,
		inputs,
		&InvokeConfig{Model: "test-model"},
	)
	require.NoError(t, err)
	require.NotNil(t, channel)

	types := make([]message.MessageType, 0)
	for msg := range channel {
		types = append(types, msg.GetType())
	}

	assert.Equal(t, []message.MessageType{
		message.AgentStarted,
		message.Text,
		message.AgentEnded,
	}, types)
	assert.Len(t, state.Capture(), 2, "input and assistant text must be kept in state")
	require.Len(t, adapter.CallInputs(), 1)
	assert.Len(t, adapter.CallInputs()[0], 1)
}
