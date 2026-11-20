package agent

import (
	"fmt"

	"github.com/smtdfc/nagare/pkgs/message"
)

type State struct {
	CurrentMessage []message.Message
	PendingMessage []message.Message
	ToolCalls      []message.ToolCallMessage
	LoopCount      int
}

func (s *State) AppendMessages(inputs []message.Message) {
	fmt.Println("Append", inputs)
	s.PendingMessage = append(s.PendingMessage, inputs...)
}

func (s *State) Capture() []message.Message {
	list := make([]message.Message, 0,
		len(s.PendingMessage)+len(s.CurrentMessage),
	)

	list = append(list, s.CurrentMessage...)
	list = append(list, s.PendingMessage...)

	return list
}

func (s *State) PrepareNextTurn() {
	s.CurrentMessage = append(s.CurrentMessage, s.PendingMessage...)
	s.PendingMessage = make([]message.Message, 0)
	s.ToolCalls = make([]message.ToolCallMessage, 0)
}

func NewState() *State {
	return &State{
		CurrentMessage: make([]message.Message, 0),
		PendingMessage: make([]message.Message, 0),
		ToolCalls:      make([]message.ToolCallMessage, 0),
		LoopCount:      0,
	}
}
