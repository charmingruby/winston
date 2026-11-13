package aikit

import (
	"context"
	"errors"
	"fmt"

	genai "github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core/api"
	"github.com/firebase/genkit/go/core/status"
	"github.com/firebase/genkit/go/genkit"

	"github.com/charmingruby/stan/internal/platform/o11y/logging"
)

var ErrInference = errors.New("inference failed")

type AIKit struct {
	Genkit *genkit.Genkit
}

func New(ctx context.Context, providers ...Provider) *AIKit {
	plugins := make([]api.Plugin, 0, len(providers))
	for _, provider := range providers {
		plugins = append(plugins, provider.Plugin())
	}

	return &AIKit{
		Genkit: genkit.Init(ctx, genkit.WithExperimental(), genkit.WithPlugins(plugins...)),
	}
}

func (a *AIKit) Infer(
	ctx context.Context,
	model genai.ModelRef,
	tools []Tool,
	opts ...genai.GenerateOption,
) (*genai.ModelResponse, error) {
	log := logging.LoggerFromContext(ctx)

	toolRefs := make([]genai.ToolRef, 0, len(tools))
	for _, tool := range tools {
		toolRefs = append(toolRefs, tool.Ref)
	}

	opts = append([]genai.GenerateOption{
		genai.WithTools(toolRefs...),
		genai.WithMaxTurns(8),
		genai.WithModel(model),
	}, opts...)

	response, err := genkit.Generate(ctx, a.Genkit, opts...)
	if err != nil {
		log.ErrorContext(ctx, "inference failed",
			"model", model.Name(),
			"status", status.Of(err),
		)

		return nil, fmt.Errorf("%w: %w", ErrInference, err)
	}

	unexpectedFinish := response.FinishReason != genai.FinishReasonStop &&
		response.FinishReason != genai.FinishReasonInterrupted
	if unexpectedFinish {
		return nil, fmt.Errorf("%w: finish reason %s", ErrInference, response.FinishReason)
	}

	return response, nil
}
