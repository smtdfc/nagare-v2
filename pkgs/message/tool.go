package message

import "github.com/google/uuid"

var (
	ToolCall MessageType = "ToolCall"
)

type ToolCallMessage struct {
	ID       string
	InvokeID string
	Type     MessageType

	CallID       string
	Name         string
	RawArguments string
}

func (t *ToolCallMessage) GetID() string {
	return t.ID
}

func (t *ToolCallMessage) GetType() MessageType {
	return t.Type
}

func (t *ToolCallMessage) GetInvokeID() string {
	return t.InvokeID
}

func (t *ToolCallMessage) SetInvokeID(s string) {
	t.InvokeID = s
}

func NewToolCallMessage(callID, name, rawArgs string) *ToolCallMessage {
	return &ToolCallMessage{
		ID:       uuid.New().String(),
		InvokeID: "",
		Type:     ToolCall,

		CallID:       callID,
		Name:         name,
		RawArguments: rawArgs,
	}
}
