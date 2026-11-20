package agent

import (
	"context"
	"sync"

	message_utils "github.com/smtdfc/nagare/core/message-utils"
	modelprovider "github.com/smtdfc/nagare/core/model-provider"
	"github.com/smtdfc/nagare/pkgs/message"
)

// fakeAdapter replies with one scripted turn per Invoke call and records the
// inputs it received. When the script is exhausted it replies with an empty
// stream.
type fakeAdapter struct {
	mu    sync.Mutex
	turns [][]message.Message
	calls [][]message.Message
	err   error
}

func (f *fakeAdapter) Invoke(ctx context.Context, inputs []message.Message, config *modelprovider.InvokeConfig) (message_utils.MessageChannel, error) {
	f.mu.Lock()
	if f.err != nil {
		err := f.err
		f.mu.Unlock()
		return nil, err
	}
	var reply []message.Message
	if len(f.turns) > 0 {
		reply = f.turns[0]
		f.turns = f.turns[1:]
	}
	f.calls = append(f.calls, append([]message.Message(nil), inputs...))
	f.mu.Unlock()

	channel := make(message_utils.MessageChannel)
	go func() {
		defer close(channel)
		for _, msg := range reply {
			select {
			case channel <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()

	return channel, nil
}

func (f *fakeAdapter) CallInputs() [][]message.Message {
	f.mu.Lock()
	defer f.mu.Unlock()

	calls := make([][]message.Message, len(f.calls))
	copy(calls, f.calls)
	return calls
}
