package message_test

import (
	"testing"

	"github.com/smtdfc/nagare/pkgs/message"
	"github.com/stretchr/testify/assert"
)

func TestReasoningMessage_GetID(t *testing.T) {
	m := message.NewReasoningMessage("t")
	assert.NotEqual(t, "", m.GetID(), "GetID() should not be empty")
}

func TestReasoningMessage_GetType(t *testing.T) {
	m := message.NewReasoningMessage("t")
	assert.Equal(t, message.Reasoning, m.GetType(), "GetType() should return Reasoning")
}

func TestReasoningMessage_GetInvokeID(t *testing.T) {
	m := message.NewReasoningMessage("t")
	assert.Equal(t, "", m.GetInvokeID(), "GetInvokeID() should be empty")
}

func TestReasoningMessage_SetInvokeID(t *testing.T) {
	m := message.NewReasoningMessage("t")
	m.SetInvokeID("t")
	assert.Equal(t, "t", m.GetInvokeID(), "SetInvokeID() should set InvokeID")
}

func TestNewReasoningMessage(t *testing.T) {
	m := message.NewReasoningMessage("thinking")
	assert.NotNil(t, m, "NewReasoningMessage() should not be nil")
	assert.Equal(t, "thinking", m.Content, "NewReasoningMessage() should set Content")
	assert.Equal(t, message.Reasoning, m.Type, "NewReasoningMessage() should have the correct type")
}
