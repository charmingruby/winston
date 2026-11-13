package aikit

import (
	"context"

	"github.com/charmingruby/winstom/internal/platform/validator"
	genai "github.com/firebase/genkit/go/ai"
)

type Tool struct {
	Ref genai.Tool
}

type ToolFn[In, Out any] func(context.Context, In) (Out, error)

func DefineTool[In, Out any](
	name, description string,
	fn ToolFn[In, Out],
) Tool {
	return Tool{
		Ref: genai.NewTool(name, description,
			func(ctx *genai.ToolContext, input In) (Out, error) {
				if err := validator.Validate(input); err != nil {
					var output Out
					return output, err
				}

				return fn(ctx, input)
			},
		),
	}
}
