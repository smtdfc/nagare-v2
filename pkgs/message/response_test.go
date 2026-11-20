package message_test

import (
	"testing"

	"github.com/smtdfc/nagare/pkgs/message"
	"github.com/stretchr/testify/assert"
)

func TestResponseFailedMessage_GetID(t *testing.T) {
	m := message.NewResponseFailedMessage("")
	assert.NotEqual(t, "", m.GetID(), "GetID() should not be empty")
}

func TestResponseFailedMessage_GetType(t *testing.T) {
	m := message.NewResponseFailedMessage("")
	assert.Equal(t, message.ResponseFailed, m.GetType(), "GetType() should return ResponseFailed")
}

func TestResponseFailedMessage_GetInvokeID(t *testing.T) {
	m := message.NewResponseFailedMessage("")
	assert.Equal(t, "", m.GetInvokeID(), "GetInvokeID() should be empty")
}

func TestResponseFailedMessage_SetInvokeID(t *testing.T) {
	m := message.NewResponseFailedMessage("")
	m.SetInvokeID("1")
	assert.Equal(t, "1", m.GetInvokeID(), "SetInvokeID() should set InvokeID")
}
func TestNewResponseFailedMessage(t *testing.T) {
	m := message.NewResponseFailedMessage("Hello")
	assert.NotNil(t, m, "NewResponseFailedMessage should not return nil")
	assert.Equal(t, "Hello", m.Error, "NewResponseFailedMessage should set Error")
	assert.Equal(t, message.ResponseFailed, m.Type, "NewResponseFailedMessage should have the correct type")
}
