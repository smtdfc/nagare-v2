package modelprovider

import (
	"context"

	messageutils "github.com/smtdfc/nagare/core/messageutils"
	"github.com/smtdfc/nagare/pkgs/message"
)

type InvokeConfig struct {
	Model string
}

type Adapter interface {
	Invoke(ctx context.Context, inputs []message.Message, config *InvokeConfig) (messageutils.MessageChannel, error)
}
