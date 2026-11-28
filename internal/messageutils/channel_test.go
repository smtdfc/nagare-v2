package message_utils

import (
	"testing"

	"github.com/smtdfc/nagare/pkgs/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageChannel_CarriesMessages(t *testing.T) {
	channel := make(MessageChannel, 1)

	channel <- message.NewTextMessage("hello", message.User)
	close(channel)

	msg, ok := <-channel
	require.True(t, ok)
	assert.Equal(t, message.Text, msg.GetType())
	assert.Equal(t, "hello", msg.(*message.TextMessage).Content)
}
