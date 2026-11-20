package modelprovider

import (
	"context"

	message_utils "github.com/smtdfc/nagare/core/message-utils"
	"github.com/smtdfc/nagare/pkgs/message"
)

type InvokeConfig struct {
	Model string
}

type Adapter interface {
	Invoke(ctx context.Context, inputs []message.Message, config *InvokeConfig) (message_utils.MessageChannel, error)
}
