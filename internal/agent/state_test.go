package agent

import (
	"testing"

	"github.com/smtdfc/nagare/pkgs/message"
	"github.com/stretchr/testify/assert"
)

func TestState_AppendMessages(t *testing.T) {
	state := NewState()
	assert.Len(t, state.CurrentMessage, 0)
	assert.Len(t, state.PendingMessage, 0)

	state.AppendMessages([]message.Message{
		message.NewAgentStartedMessage(),
	})

	assert.Len(t, state.CurrentMessage, 0)
	assert.Len(t, state.PendingMessage, 1)
}

func TestState_Capture(t *testing.T) {
	tests := []struct {
		name            string
		currentMessages []message.Message
		pendingMessages []message.Message
		expected        []message.MessageType
	}{
		{
			name:     "empty state",
			expected: []message.MessageType{},
		},
		{
			name: "only pending messages",
			pendingMessages: []message.Message{
				message.NewAgentStartedMessage(),
			},
			expected: []message.MessageType{
				message.AgentStarted,
			},
		},
		{
			name: "only current messages",
			currentMessages: []message.Message{
				message.NewAgentStartedMessage(),
			},
			expected: []message.MessageType{
				message.AgentStarted,
			},
		},
		{
			name: "both current and pending messages",
			currentMessages: []message.Message{
				message.NewAgentStartedMessage(),
			},
			pendingMessages: []message.Message{
				message.NewAgentStartedMessage(),
				message.NewAgentStartedMessage(),
			},
			expected: []message.MessageType{
				message.AgentStarted,
				message.AgentStarted,
				message.AgentStarted,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := NewState()
			state.CurrentMessage = tt.currentMessages
			state.PendingMessage = tt.pendingMessages

			got := state.Capture()

			assert.Len(t, got, len(tt.expected))

			for i, expected := range tt.expected {
				assert.Equal(t, expected, got[i].GetType())
			}
		})
	}
}

func TestState_Capture_IndependentSnapshot(t *testing.T) {
	state := NewState()

	state.AppendMessages([]message.Message{
		message.NewAgentStartedMessage(),
	})

	snapshot := state.Capture()

	state.AppendMessages([]message.Message{
		message.NewAgentStartedMessage(),
	})

	assert.Len(t, snapshot, 1)
	assert.Len(t, state.PendingMessage, 2)
}

func TestNewState(t *testing.T) {
	state := NewState()
	assert.NotNil(t, state)
	assert.NotNil(t, state.CurrentMessage)
	assert.NotNil(t, state.PendingMessage)
	assert.Equal(t, 0, state.LoopCount)
	assert.Len(t, state.CurrentMessage, 0)
	assert.Len(t, state.PendingMessage, 0)
}

func TestState_PrepareNextTurn(t *testing.T) {
	state := NewState()
	state.AppendMessages([]message.Message{
		message.NewTextMessage("answer", message.Agent),
	})
	state.ToolCalls = append(state.ToolCalls, *message.NewToolCallMessage("call-1", "lookup", "{}"))

	state.PrepareNextTurn()

	assert.Len(t, state.CurrentMessage, 1)
	assert.Len(t, state.PendingMessage, 0)
	assert.Len(t, state.ToolCalls, 0)

	state.AppendMessages([]message.Message{message.NewAgentStartedMessage()})

	assert.Len(t, state.CurrentMessage, 1, "new pending messages must not overwrite committed history")
	assert.Len(t, state.PendingMessage, 1)
}
