package message

import "github.com/google/uuid"

var (
	ResponseFailed MessageType = "ResponseFailed"
)

type ResponseFailedMessage struct {
	ID       string      `json:"id"`
	InvokeID string      `json:"invoke_id"`
	Type     MessageType `json:"type"`
	Error    string      `json:"error"`
}

func (r *ResponseFailedMessage) GetID() string {
	return r.ID
}

func (r *ResponseFailedMessage) GetType() MessageType {
	return r.Type
}

func (r *ResponseFailedMessage) GetInvokeID() string {
	return r.InvokeID
}

func (r *ResponseFailedMessage) SetInvokeID(s string) {
	r.InvokeID = s
}

func NewResponseFailedMessage(err string) *ResponseFailedMessage {
	return &ResponseFailedMessage{
		ID:       uuid.New().String(),
		InvokeID: "",
		Type:     ResponseFailed,
		Error:    err,
	}
}
