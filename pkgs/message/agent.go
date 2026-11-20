package message

import "github.com/google/uuid"

var (
	AgentStarted MessageType = "AgentStarted"
	AgentEnded   MessageType = "AgentEnded"
)

type AgentStartedMessage struct {
	ID       string
	InvokeID string
	Type     MessageType
}

// GetID implements [Message].
func (a *AgentStartedMessage) GetID() string {
	return a.ID
}

// GetInvokeID implements [Message].
func (a *AgentStartedMessage) GetInvokeID() string {
	return a.InvokeID
}

// GetType implements [Message].
func (a *AgentStartedMessage) GetType() MessageType {
	return a.Type
}

// SetInvokeID implements [Message].
func (a *AgentStartedMessage) SetInvokeID(id string) {
	a.InvokeID = id
}

func NewAgentStartedMessage() *AgentStartedMessage {
	return &AgentStartedMessage{
		ID:       uuid.New().String(),
		Type:     AgentStarted,
		InvokeID: "",
	}
}

type AgentEndedMessage struct {
	ID          string
	InvokeID    string
	Type        MessageType
	IsSuccess   bool
	Error       string
	ExecuteTime float64
}

// GetID implements [Message].
func (a *AgentEndedMessage) GetID() string {
	return a.ID
}

// GetInvokeID implements [Message].
func (a *AgentEndedMessage) GetInvokeID() string {
	return a.InvokeID
}

// GetType implements [Message].
func (a *AgentEndedMessage) GetType() MessageType {
	return a.Type
}

// SetInvokeID implements [Message].
func (a *AgentEndedMessage) SetInvokeID(id string) {
	a.InvokeID = id
}

func NewAgentEndedMessage(
	isSuccess bool,
	executeTime float64,
	err string,
) *AgentEndedMessage {
	return &AgentEndedMessage{
		ID:          uuid.New().String(),
		InvokeID:    "",
		Type:        AgentEnded,
		IsSuccess:   isSuccess,
		ExecuteTime: executeTime,
		Error:       err,
	}
}
