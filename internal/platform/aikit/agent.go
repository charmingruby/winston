package aikit

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/charmingruby/stan/internal/platform/o11y/logging"
	genai "github.com/firebase/genkit/go/ai"
	aix "github.com/firebase/genkit/go/ai/exp"
	"github.com/firebase/genkit/go/genkit"
	genkitx "github.com/firebase/genkit/go/genkit/exp"
)

type Agent struct {
	def  *aix.Agent[struct{}]
	name string
}

func DefineAgent(
	ak *AIKit,
	name string,
	model genai.ModelRef,
	tools []Tool,
	prompt aix.InlinePrompt,
	opts ...aix.AgentOption[struct{}],
) (*Agent, error) {
	refs := make([]genai.ToolRef, 0, len(tools))
	for _, tool := range tools {
		genkit.RegisterAction(ak.Genkit, tool.Ref)
		refs = append(refs, tool.Ref)
	}

	cfg := aix.InlinePrompt{
		genai.WithTools(refs...),
		genai.WithMaxTurns(8),
		genai.WithModel(model),
		genai.WithUse(
			toolPolicy(),
			agentObserver(name),
		),
	}

	cfg = append(cfg, prompt...)

	def := genkitx.DefineAgent[struct{}](ak.Genkit, name, cfg, opts...)

	return &Agent{
		def:  def,
		name: name,
	}, nil
}

func (a *Agent) Run(
	ctx context.Context,
	message string,
	opts ...aix.InvocationOption[struct{}],
) (string, error) {
	log := logging.LoggerFromContext(ctx)

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	out, err := a.def.RunText(ctx, message, opts...)
	if err != nil {
		log.WarnContext(ctx, "agent invocation failed",
			"agent", a.name,
		)

		return "", fmt.Errorf("%w: %w", ErrInference, err)
	}

	log.InfoContext(ctx, "agent execution finished", "agent", a.name, "outcome", out.FinishReason)

	if out.Usage != nil {
		log.InfoContext(ctx, "agent model usage",
			"agent", a.name,
			"input_tokens", out.Usage.InputTokens,
			"output_tokens", out.Usage.OutputTokens,
		)
	}

	if out.Error != nil {
		return "", fmt.Errorf("%w: %w", ErrInference, out.Error)
	}

	unexpectedFinish := out.FinishReason != aix.AgentFinishReasonStop || out.Message == nil

	if unexpectedFinish {
		return "", fmt.Errorf("%w: agent finish reason %s", ErrInference, out.FinishReason)
	}

	text := out.Message.Text()
	if text == "" {
		return "", fmt.Errorf("%w: empty agent response", ErrInference)
	}

	return text, nil
}

func (a *Agent) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	genkit.Handler(a.def).ServeHTTP(w, r)
}

func (a *Agent) Ref() aix.AgentRef {
	return a.def.Ref()
}

func toolPolicy() genai.Middleware {
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

func agentObserver(name string) genai.Middleware {
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
