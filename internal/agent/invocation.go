package agent

import (
	errors "github.com/smtdfc/nagare/core/errors"
	"github.com/smtdfc/nagare/core/execution"
	messageutils "github.com/smtdfc/nagare/core/messageutils"
	modelprovider "github.com/smtdfc/nagare/core/modelprovider"
	"github.com/smtdfc/nagare/pkgs/message"
)

type InvokeConfig struct {
	Model string
}

func (a *Agent) Invoke(ctx *execution.Context, adapter modelprovider.Adapter, inputs []message.Message, config *InvokeConfig) (messageutils.MessageChannel, error) {
	if a.State == nil {
		return nil, errors.ErrAgentStateEmpty
	}

	if adapter == nil {
		return nil, errors.ErrAdapterNotFound
	}

	a.State.AppendMessages(inputs)
	outputChannel := make(chan message.Message)
	go func() {
		err := NewExecutor(ctx, a.logger).
			WithChannel(outputChannel).
			WithAdapter(adapter).
			WithState(a.State).
			WithConfig(config).
			Start()
		if err != nil {
			a.logger.Error("agent execution failed", "error", err)
		}
	}()

	return outputChannel, nil
}
