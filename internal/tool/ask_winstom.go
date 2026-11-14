package tool

import (
	"context"
	"fmt"

	"github.com/charmingruby/erebor/aikit"
)

type AskWinstonInput struct {
	Message string `json:"message" validate:"required" jsonschema:"description=Overwatch question or task for winston. Include all relevant context because each call starts a new conversation."`
}

type AskWinstonOutput struct {
	Answer string `json:"answer" jsonschema:"description=winston's response after consulting its specialist agents and team tools as needed."`
}

func AskWinston(run func(context.Context, string) (string, error)) aikit.Tool {
	return aikit.DefineTool(
		"ask_winston",
		"Asks winston for Overwatch team composition, strategy, or gameplay analysis. "+
			"winston can consult the team roster and delegate to analyst, strategist, and reviewer agents. "+
			"Each call is independent; include relevant context from previous messages.",
		func(ctx context.Context, input AskWinstonInput) (AskWinstonOutput, error) {
			if input.Message == "" {
				return AskWinstonOutput{}, fmt.Errorf("message must not be blank")
			}

			answer, err := run(ctx, input.Message)
			if err != nil {
				return AskWinstonOutput{}, fmt.Errorf("ask winston: %w", err)
			}

			return AskWinstonOutput{Answer: answer}, nil
		},
	)
}
