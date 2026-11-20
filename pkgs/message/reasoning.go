package message

import "github.com/google/uuid"

var (
	Reasoning MessageType = "Reasoning"
)

type ReasoningMessage struct {
	ID       string
	InvokeID string
	Type     MessageType
	Content  string
}

func (r *ReasoningMessage) GetID() string {
	return r.ID
}

func (r *ReasoningMessage) GetType() MessageType {
	return r.Type
}

func (r *ReasoningMessage) GetInvokeID() string {
	return r.InvokeID
}

func (r *ReasoningMessage) SetInvokeID(s string) {
	r.InvokeID = s
}

func NewReasoningMessage(content string) *ReasoningMessage {
	return &ReasoningMessage{
		ID:       uuid.New().String(),
		InvokeID: "",
		Type:     Reasoning,
		Content:  content,
	}
}
