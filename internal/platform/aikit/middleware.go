package aikit

import (
	"context"
	"fmt"

	"github.com/charmingruby/winstom/internal/platform/o11y/logging"
	genai "github.com/firebase/genkit/go/ai"
)

func toolPolicyMiddleware() genai.Middleware {
	return genai.MiddlewareFunc(func(context.Context) (*genai.Hooks, error) {
		return &genai.Hooks{
			WrapModel: func(
				ctx context.Context,
				params *genai.ModelParams,
				next genai.ModelNext,
			) (*genai.ModelResponse, error) {
				response, err := next(ctx, params)
				if err != nil {
					return nil, err
				}

				requests := response.ToolRequests()
				if len(requests) > 1 {
					return nil, fmt.Errorf("%w: only one tool call per turn is allowed", ErrInference)
				}

				allowed := make(map[string]struct{}, len(params.Request.Tools))
				for _, tool := range params.Request.Tools {
					allowed[tool.Name] = struct{}{}
				}

				for _, part := range requests {
					if _, ok := allowed[part.ToolRequest.Name]; !ok {
						return nil, fmt.Errorf("%w: unavailable tool %q", ErrInference, part.ToolRequest.Name)
					}
				}

				return response, nil
			},
		}, nil
	})
}

func agentObserverMiddleware(name string) genai.Middleware {
	return genai.MiddlewareFunc(func(context.Context) (*genai.Hooks, error) {
		return &genai.Hooks{
			WrapTool: func(
				ctx context.Context,
				params *genai.ToolParams,
				next genai.ToolNext,
			) (*genai.MultipartToolResponse, error) {
				result, err := next(ctx, params)

				logging.LoggerFromContext(ctx).InfoContext(ctx, "agent tool finished",
					"agent", name,
					"tool", params.Request.Name,
					"success", err == nil,
				)

				return result, err
			},
		}, nil
	})
}
