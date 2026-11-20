package agent

import (
	error_utils "github.com/smtdfc/nagare/core/error-utils"
	"github.com/smtdfc/nagare/core/execution"
	message_utils "github.com/smtdfc/nagare/core/message-utils"
	modelprovider "github.com/smtdfc/nagare/core/model-provider"
	"github.com/smtdfc/nagare/pkgs/message"
)

type InvokeConfig struct {
	Model string
}

func (a *Agent) Invoke(ctx *execution.Context, adapter modelprovider.Adapter, inputs []message.Message, config *InvokeConfig) (message_utils.MessageChannel, error) {
	if a.State == nil {
		return nil, error_utils.ErrAgentStateEmpty
	}

	if adapter == nil {
		return nil, error_utils.ErrAdapterNotFound
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
