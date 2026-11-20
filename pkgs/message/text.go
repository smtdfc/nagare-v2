package message

import "github.com/google/uuid"

var (
	Text MessageType = "text"
)

type Role string

var (
	User   Role = "user"
	Agent  Role = "agent"
	System Role = "system"
)

type TextMessage struct {
	ID       string
	InvokeID string
	Type     MessageType
	Content  string
	Role     Role
}

// GetID implements [Message].
func (t *TextMessage) GetID() string {
	return t.ID
}

// GetInvokeID implements [Message].
func (t *TextMessage) GetInvokeID() string {
	return t.InvokeID
}

// GetType implements [Message].
func (t *TextMessage) GetType() MessageType {
	return t.Type
}

// SetInvokeID implements [Message].
func (t *TextMessage) SetInvokeID(id string) {
	t.InvokeID = id
}

func NewTextMessage(content string, role Role) *TextMessage {
	return &TextMessage{
		ID:       uuid.New().String(),
		InvokeID: "",
		Type:     Text,
		Content:  content,
		Role:     role,
	}
}
