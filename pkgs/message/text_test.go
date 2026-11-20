package message_test

import (
	"testing"

	"github.com/smtdfc/nagare/pkgs/message"
	"github.com/stretchr/testify/assert"
)

func TestTextMessage_GetID(t *testing.T) {
	m := message.NewTextMessage("hello", message.Agent)
	assert.NotEqual(t, "", m.GetID(), "GetID() should not be empty")
}

func TestTextMessage_GetInvokeID(t *testing.T) {
	m := message.NewTextMessage("hello", message.Agent)
	assert.Equal(t, "", m.GetInvokeID(), "GetInvokeID() should be empty")
}

func TestTextMessage_SetInvokeID(t *testing.T) {
	m := message.NewTextMessage("hello", message.Agent)
	m.SetInvokeID("hello")
	assert.Equal(t, m.GetInvokeID(), "hello", "SetInvokeID() should be equal")
}

func TestNewTextMessage(t *testing.T) {
	m := message.NewTextMessage("hello", message.Agent)
	assert.NotNil(t, m, "NewTextMessage() should not be nil")
	assert.Equal(t, message.Text, m.Type, "NewTextMessage() should have the correct type")
}
