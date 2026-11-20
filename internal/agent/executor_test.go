package agent

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/smtdfc/nagare/core/execution"
	message_utils "github.com/smtdfc/nagare/core/message-utils"
	modelprovider "github.com/smtdfc/nagare/core/model-provider"
	"github.com/smtdfc/nagare/pkgs/logger"
	"github.com/smtdfc/nagare/pkgs/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testingAdapter struct {
	testTypes []string
}

func (t testingAdapter) Invoke(ctx context.Context, inputs []message.Message, config *modelprovider.InvokeConfig) (message_utils.MessageChannel, error) {
	channel := make(message_utils.MessageChannel)

	go func() {
		defer close(channel)
		if slices.Contains(t.testTypes, "text") {
			for i := range 3 {
				channel <- message.NewTextMessage(fmt.Sprintf("hello world %d", i), message.Agent)
			}
		}

		if slices.Contains(t.testTypes, "tool") {
			channel <- message.NewToolCallMessage("1234", "abc", "{}")
		}
	}()

	return channel, nil
}

func TestExecutor_WithAdapter(t *testing.T) {
	var ctx = execution.NewContext(context.Background())
	e := NewExecutor(ctx, logger.NewLogger())
	assert.Nil(t, e.adapter)
	e.WithAdapter(&testingAdapter{})
	assert.NotNil(t, e.adapter)
}

func TestExecutor_WithState(t *testing.T) {
	var ctx = execution.NewContext(context.Background())
	e := NewExecutor(ctx, logger.NewLogger())
	state := NewState()
	assert.Nil(t, e.state)
	e.WithState(state)
	assert.NotNil(t, e.state)
}

func TestExecutor_WithChannel(t *testing.T) {
	var ctx = execution.NewContext(context.Background())
	channel := make(message_utils.MessageChannel)
	e := NewExecutor(ctx, logger.NewLogger())
	assert.Nil(t, e.outputChan)
	e.WithChannel(channel)
	assert.NotNil(t, e.outputChan)
}

func TestExecutor_WithConfig(t *testing.T) {
	var ctx = execution.NewContext(context.Background())
	config := &InvokeConfig{}
	e := NewExecutor(ctx, logger.NewLogger())
	assert.Nil(t, e.config)
	e.WithConfig(config)
	assert.NotNil(t, e.config)
}

func TestExecutor_Start_Stream(t *testing.T) {
	var ctx = execution.NewContext(context.Background())
	state := NewState()
	adapter := &testingAdapter{
		testTypes: []string{"text"},
	}
	outputChan := make(message_utils.MessageChannel)
	e := NewExecutor(ctx, logger.NewLogger())
	e.WithAdapter(adapter)
	e.WithState(state)
	e.WithConfig(&InvokeConfig{})
	e.WithChannel(outputChan)
	go e.Start()
	messages := make([]message.Message, 0)
	for msg := range outputChan {
		messages = append(messages, msg)
	}

	fmt.Println(messages)
	assert.Len(t, messages, 5)
	assert.Equal(t, message.Text, messages[1].GetType())
	assert.Equal(t, message.Text, messages[3].GetType())
	assert.Equal(t, message.AgentEnded, messages[4].GetType())
	fmt.Println("Capture: ", state.Capture())
	assert.Len(t, state.Capture(), 1)
}

func TestNewExecutor(t *testing.T) {
	var ctx = execution.NewContext(context.Background())
	e := NewExecutor(ctx, logger.NewLogger())
	assert.NotNil(t, e)
	assert.NotNil(t, e.ctx)
	assert.NotNil(t, e.logger)
}

func TestNewExecutor_DefaultLogger(t *testing.T) {
	e := NewExecutor(execution.NewContext(context.Background()), nil)
	assert.NotNil(t, e.logger)
}

func TestExecutor_Start_ContinuesAfterToolCall(t *testing.T) {
	state := NewState()
	adapter := &fakeAdapter{
		turns: [][]message.Message{
			{
				message.NewTextMessage("thinking", message.Agent),
				message.NewToolCallMessage("call-1", "lookup", "{}"),
			},
			{message.NewTextMessage("final", message.Agent)},
		},
	}
	outputChan := make(message_utils.MessageChannel)
	e := NewExecutor(execution.NewContext(context.Background()), logger.NewLogger()).
		WithAdapter(adapter).
		WithState(state).
		WithConfig(&InvokeConfig{Model: "test-model"}).
		WithChannel(outputChan)

	go e.Start()
	for range outputChan {
	}

	calls := adapter.CallInputs()
	require.Len(t, calls, 2, "a tool call must trigger a second model turn")
	assert.Empty(t, calls[0])
	require.Len(t, calls[1], 2)
	assert.Equal(t, message.Text, calls[1][0].GetType())
	assert.Equal(t, message.ToolCall, calls[1][1].GetType())
	assert.Len(t, state.Capture(), 3)
}

func TestExecutor_Start_AdapterError(t *testing.T) {
	state := NewState()
	adapter := &fakeAdapter{err: errors.New("boom")}
	outputChan := make(message_utils.MessageChannel)
	e := NewExecutor(execution.NewContext(context.Background()), logger.NewLogger()).
		WithAdapter(adapter).
		WithState(state).
		WithConfig(&InvokeConfig{Model: "test-model"}).
		WithChannel(outputChan)

	go e.Start()

	messages := make([]message.Message, 0)
	for msg := range outputChan {
		messages = append(messages, msg)
	}

	require.Len(t, messages, 2)
	assert.Equal(t, message.AgentStarted, messages[0].GetType())
	ended, ok := messages[1].(*message.AgentEndedMessage)
	require.True(t, ok)
	assert.False(t, ended.IsSuccess)
	assert.Equal(t, "boom", ended.Error)
	assert.Empty(t, state.Capture())
}
