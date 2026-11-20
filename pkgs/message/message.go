package message

type MessageType string

type Message interface {
	GetID() string
	GetType() MessageType
	GetInvokeID() string
	SetInvokeID(string)
}
