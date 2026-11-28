package agent

import (
	"strings"
	"time"

	errors "github.com/smtdfc/nagare/core/errors"
	"github.com/smtdfc/nagare/core/execution"
	messageutils "github.com/smtdfc/nagare/core/messageutils"
	modelprovider "github.com/smtdfc/nagare/core/modelprovider"
	"github.com/smtdfc/nagare/pkgs/logger"
	"github.com/smtdfc/nagare/pkgs/message"
)

const maxIteration = 20

type Executor struct {
	ctx         *execution.Context
	outputChan  messageutils.MessageChannel
	adapter     modelprovider.Adapter
	config      *InvokeConfig
	state       *State
	logger      *logger.Logger
	textBuilder strings.Builder
}

func (e *Executor) flushTextMessage() {
	if e.textBuilder.Len() > 0 {
		e.state.AppendMessages([]message.Message{
			message.NewTextMessage(e.textBuilder.String(), message.Agent),
		})
		e.textBuilder.Reset()
	}
}

func (e *Executor) handleMessage(stream messageutils.MessageChannel) {
	defer e.flushTextMessage()

	for {
		select {
		case <-e.ctx.Done():
			return

		case chunk, ok := <-stream:
			if !ok {
				return
			}

			switch v := chunk.(type) {
			case *message.TextMessage:
				e.textBuilder.WriteString(v.Content)
			case *message.ToolCallMessage:
				e.flushTextMessage()
				e.state.AppendMessages([]message.Message{
					chunk,
				})
				e.state.ToolCalls = append(e.state.ToolCalls, *v)
			default:
				e.flushTextMessage()
				e.logger.Debug("message not persisted in state", "type", chunk.GetType())
			}
			select {
			case e.outputChan <- chunk:
			case <-e.ctx.Done():
				return
			}
		}
	}
}

func (e *Executor) startTurn() error {
	if e.state.loopCount >= maxIteration {
		return errors.ErrAgentLoopMaxIterations
	}

	stream, err := e.adapter.Invoke(e.ctx, e.state.Capture(), &modelprovider.InvokeConfig{
		Model: e.config.Model,
	})
	if err != nil {
		return err
	}

	e.handleMessage(stream)
	e.state.loopCount++
	return nil
}

func (e *Executor) nextTurn() bool {
	return len(e.state.ToolCalls) > 0
}

func (e *Executor) Start() error {
	startTime := time.Now()
	e.outputChan <- message.NewAgentStartedMessage()

	defer close(e.outputChan)

	for {
		if err := e.ctx.Err(); err != nil {
			e.logger.Info("agent canceled", "reason", err)
			return err
		}

		if err := e.startTurn(); err != nil {
			e.logger.Error("agent turn failed", "error", err, "loops", e.state.loopCount)

			endTime := time.Now()
			e.outputChan <- message.NewAgentEndedMessage(
				false,
				endTime.Sub(startTime).Seconds(),
				err.Error(),
			)
			return err
		}

		if err := e.ctx.Err(); err != nil {
			e.logger.Info("agent canceled", "reason", err)
			return err
		}

		if !e.nextTurn() {
			break
		}

		e.executeTools()
		e.state.PrepareNextTurn()
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds()
	e.logger.Info("agent finished", "loops", e.state.loopCount, "duration", duration)

	e.outputChan <- message.NewAgentEndedMessage(
		true,
		duration,
		"",
	)

	return nil
}

func (e *Executor) executeTools() {
	for _, call := range e.state.ToolCalls {
		select {
		case <-e.ctx.Done():
			e.logger.Info("tool execution canceled", "reason", e.ctx.Err())
			return
		default:
		}

		e.logger.Warn("tool execution is not implemented", "tool", call.Name, "call_id", call.CallID)
	}
}

func (e *Executor) WithChannel(channel messageutils.MessageChannel) *Executor {
	e.outputChan = channel
	return e
}

func (e *Executor) WithAdapter(adapter modelprovider.Adapter) *Executor {
	e.adapter = adapter
	return e
}

func (e *Executor) WithConfig(config *InvokeConfig) *Executor {
	e.config = config
	return e
}

func (e *Executor) WithState(s *State) *Executor {
	e.state = s
	return e
}

func NewExecutor(ctx *execution.Context, l *logger.Logger) *Executor {
	if l == nil {
		l = logger.NewLogger()
	}

	return &Executor{
		ctx:    ctx,
		logger: l,
	}
}
