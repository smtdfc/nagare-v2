package adapters

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openai/openai-go/v3/responses"
	errors "github.com/smtdfc/nagare/core/errors"
	modelprovider "github.com/smtdfc/nagare/core/modelprovider"
	"github.com/smtdfc/nagare/pkgs/logger"
	"github.com/smtdfc/nagare/pkgs/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOpenAIAdapter(t *testing.T) {
	logger := logger.NewLogger()
	adp := NewOpenAIAdapter("", "", logger)
	assert.NotNil(t, adp)
	assert.NotNil(t, adp.client)
	assert.NotNil(t, adp.logger)
	assert.Equal(t, logger, adp.logger)
}

func TestOpenAIAdapter_toProviderMessage(t *testing.T) {
	adp := NewOpenAIAdapter("", "", logger.NewLogger())

	t.Run("user text becomes input message", func(t *testing.T) {
		out, err := adp.toProviderMessage(message.NewTextMessage("Hello", message.User))

		require.NoError(t, err)
		require.NotNil(t, out.OfInputMessage)
		assert.Equal(t, "user", out.OfInputMessage.Role)
		require.Len(t, out.OfInputMessage.Content, 1)
		assert.Equal(t, "Hello", out.OfInputMessage.Content[0].OfInputText.Text)
	})

	t.Run("system text becomes input message", func(t *testing.T) {
		out, err := adp.toProviderMessage(message.NewTextMessage("Rules", message.System))

		require.NoError(t, err)
		require.NotNil(t, out.OfInputMessage)
		assert.Equal(t, "system", out.OfInputMessage.Role)
	})

	t.Run("agent text becomes output message", func(t *testing.T) {
		out, err := adp.toProviderMessage(message.NewTextMessage("Answer", message.Agent))

		require.NoError(t, err)
		require.NotNil(t, out.OfOutputMessage)
		require.Len(t, out.OfOutputMessage.Content, 1)
		assert.Equal(t, "Answer", out.OfOutputMessage.Content[0].OfOutputText.Text)
	})

	t.Run("tool call becomes function call", func(t *testing.T) {
		out, err := adp.toProviderMessage(message.NewToolCallMessage("call-1", "lookup", `{"q":"x"}`))

		require.NoError(t, err)
		require.NotNil(t, out.OfFunctionCall)
		assert.Equal(t, "call-1", out.OfFunctionCall.CallID)
		assert.Equal(t, "lookup", out.OfFunctionCall.Name)
		assert.Equal(t, `{"q":"x"}`, out.OfFunctionCall.Arguments)
	})

	t.Run("unsupported message becomes empty item", func(t *testing.T) {
		out, err := adp.toProviderMessage(message.NewAgentStartedMessage())

		require.NoError(t, err)
		assert.Equal(t, responses.ResponseInputItemUnionParam{}, out)
	})
}

func TestOpenAIAdapter_Invoke_StreamsEvents(t *testing.T) {
	events := []string{
		`{"type":"response.output_text.delta","delta":"Hel"}`,
		`{"type":"response.output_text.delta","delta":"lo"}`,
		`{"type":"response.output_text.delta","delta":""}`,
		`{"type":"response.reasoning_text.delta","delta":"thinking"}`,
		`{"type":"response.output_item.done","item":{"type":"function_call","call_id":"call-1","name":"lookup","arguments":"{}"}}`,
	}
	server := newSSEServer(t, http.StatusOK, events)
	adp := NewOpenAIAdapter(server.URL, "test-key", logger.NewLogger())

	channel, err := adp.Invoke(
		context.Background(),
		[]message.Message{message.NewTextMessage("question", message.User)},
		&modelprovider.InvokeConfig{Model: "test-model"},
	)
	require.NoError(t, err)
	require.NotNil(t, channel)

	got := make([]message.Message, 0)
	for msg := range channel {
		got = append(got, msg)
	}

	require.Len(t, got, 4)
	assert.Equal(t, message.Text, got[0].GetType())
	assert.Equal(t, "Hel", got[0].(*message.TextMessage).Content)
	assert.Equal(t, message.Text, got[1].GetType())
	assert.Equal(t, "lo", got[1].(*message.TextMessage).Content)
	assert.Equal(t, message.Reasoning, got[2].GetType())
	assert.Equal(t, "thinking", got[2].(*message.ReasoningMessage).Content)
	assert.Equal(t, message.ToolCall, got[3].GetType())
	toolCall := got[3].(*message.ToolCallMessage)
	assert.Equal(t, "call-1", toolCall.CallID)
	assert.Equal(t, "lookup", toolCall.Name)
	assert.Equal(t, "{}", toolCall.RawArguments)
}

func TestOpenAIAdapter_Invoke_ResponseFailedEvent(t *testing.T) {
	events := []string{
		`{"type":"response.failed","response":{"error":{"code":"server_error","message":"provider failed"}}}`,
	}
	server := newSSEServer(t, http.StatusOK, events)
	adp := NewOpenAIAdapter(server.URL, "test-key", logger.NewLogger())

	channel, err := adp.Invoke(
		context.Background(),
		nil,
		&modelprovider.InvokeConfig{Model: "test-model"},
	)
	require.NoError(t, err)

	got := make([]message.Message, 0)
	for msg := range channel {
		got = append(got, msg)
	}

	require.Len(t, got, 1)
	assert.Equal(t, message.ResponseFailed, got[0].GetType())
	assert.Equal(t, "provider failed", got[0].(*message.ResponseFailedMessage).Error)
}

func TestOpenAIAdapter_Invoke_NotFoundError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"message":"model not found"}}`))
	}))
	t.Cleanup(server.Close)

	adp := NewOpenAIAdapter(server.URL, "test-key", logger.NewLogger())
	channel, err := adp.Invoke(
		context.Background(),
		nil,
		&modelprovider.InvokeConfig{Model: "missing-model"},
	)
	require.NoError(t, err)

	got := make([]message.Message, 0)
	for msg := range channel {
		got = append(got, msg)
	}

	require.Len(t, got, 1)
	assert.Equal(t, message.ResponseFailed, got[0].GetType())
	assert.Equal(t, errors.ErrModelNotFound.Error(), got[0].(*message.ResponseFailedMessage).Error)
}

func newSSEServer(t *testing.T, status int, events []string) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(status)
		flusher, _ := w.(http.Flusher)

		for _, event := range events {
			_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
			if flusher != nil {
				flusher.Flush()
			}
		}
	}))
	t.Cleanup(server.Close)

	return server
}
