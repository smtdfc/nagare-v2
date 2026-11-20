package adapters

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
	message_utils "github.com/smtdfc/nagare/core/message-utils"
	modelprovider "github.com/smtdfc/nagare/core/model-provider"
	"github.com/smtdfc/nagare/pkgs/logger"
	"github.com/smtdfc/nagare/pkgs/message"
)

type OpenAIAdapter struct {
	logger *logger.Logger
	client *openai.Client
}

func (o *OpenAIAdapter) toProviderMessage(msg message.Message) (responses.ResponseInputItemUnionParam, error) {
	switch t := msg.(type) {
	case *message.TextMessage:
		item := &responses.ResponseInputItemMessageParam{
			Type: "message",
			Content: responses.ResponseInputMessageContentListParam{
				responses.ResponseInputContentUnionParam{
					OfInputText: &responses.ResponseInputTextParam{
						Text: t.Content,
					},
				},
			},
		}

		switch t.Role {
		case message.User:
			item.Role = "user"
			return responses.ResponseInputItemUnionParam{
				OfInputMessage: item,
			}, nil

		case message.System:
			item.Role = "system"
			return responses.ResponseInputItemUnionParam{
				OfInputMessage: item,
			}, nil
		case message.Agent:
			return responses.ResponseInputItemUnionParam{
				OfOutputMessage: &responses.ResponseOutputMessageParam{
					Content: []responses.ResponseOutputMessageContentUnionParam{
						{
							OfOutputText: &responses.ResponseOutputTextParam{
								Text: t.Content,
							},
						},
					},
				},
			}, nil
		}

	case *message.ToolCallMessage:
		return responses.ResponseInputItemUnionParam{
			OfFunctionCall: &responses.ResponseFunctionToolCallParam{
				CallID:    t.CallID,
				Name:      t.Name,
				Arguments: t.RawArguments,
			},
		}, nil
	}

	return responses.ResponseInputItemUnionParam{}, nil
}

// Invoke implements [modelprovider.Adapter].
func (o *OpenAIAdapter) Invoke(ctx context.Context, inputs []message.Message, config *modelprovider.InvokeConfig) (message_utils.MessageChannel, error) {
	outputChannel := make(chan message.Message)
	messages := make([]responses.ResponseInputItemUnionParam, 0, len(inputs))
	for _, msg := range inputs {
		t, err := o.toProviderMessage(msg)
		if err != nil {
			return nil, err
		}

		messages = append(messages, t)
	}

	go (func() {
		defer close(outputChannel)
		stream := o.client.Responses.NewStreaming(ctx, responses.ResponseNewParams{
			Model: config.Model,
			Input: responses.ResponseNewParamsInputUnion{
				OfInputItemList: messages,
			},
			//Tools:       listTool,
			Temperature: param.NewOpt(0.1),
			TopP:        param.NewOpt(0.9),
		})

		for stream.Next() {
			event := stream.Current()
			switch event.Type {
			case "response.failed":
				resp := event.AsResponseFailed().Response
				err := resp.Error
				// usage := resp.Usage
				// cb(&messages.ResponseStatsMessage{
				// 	InputTokens:     usage.InputTokens,
				// 	OutputTokens:    usage.OutputTokens,
				// 	ReasoningTokens: usage.OutputTokensDetails.ReasoningTokens,
				// 	TotalTokens:     usage.TotalTokens,
				// })
				outputChannel <- message.NewResponseFailedMessage(err.Message)
			case "response.output_item.done":
				if event.AsResponseOutputItemAdded().Item.Type == "function_call" {
					item := event.AsResponseOutputItemAdded().Item
					outputChannel <- message.NewToolCallMessage(
						item.CallID,
						item.Name,
						item.Arguments.OfString,
					)
				}

			case "response.output_text.delta":
				if event.Delta != "" {
					outputChannel <- message.NewTextMessage(
						event.Delta,
						message.Agent,
					)
				}

			case "response.reasoning_text.delta":
				if event.Delta != "" {
					outputChannel <- message.NewReasoningMessage(
						event.Delta,
					)
				}
			}
		}

		if err := stream.Err(); err != nil {
			if errors.Is(ctx.Err(), context.Canceled) || errors.Is(err, context.Canceled) {
				return
			}
			o.logger.Error("stream error", "error", err)
			if strings.Contains(err.Error(), "404") {
				outputChannel <- message.NewResponseFailedMessage(
					fmt.Sprintf("Model %s not found", config.Model),
				)
			} else if strings.Contains(err.Error(), "429") {
				outputChannel <- message.NewResponseFailedMessage(
					fmt.Sprintf("Quota exceed: %s", err.Error()),
				)
			} else if strings.Contains(err.Error(), "402") {
				outputChannel <- message.NewResponseFailedMessage(
					fmt.Sprintf("Payment required: %s", err.Error()),
				)
			} else {
				outputChannel <- message.NewResponseFailedMessage(
					err.Error(),
				)
			}
		}
	})()

	return outputChannel, nil
}

func NewOpenAIAdapter(
	baseURL string,
	apiKey string,
	logger *logger.Logger,
) *OpenAIAdapter {
	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseURL))
	return &OpenAIAdapter{
		client: &client,
		logger: logger,
	}
}
