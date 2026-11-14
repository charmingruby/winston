package workflow

import (
	"context"
	"fmt"

	"github.com/charmingruby/erebor/aikit"
	"github.com/charmingruby/erebor/aikit/provider"
	"github.com/charmingruby/winston/config"
	"github.com/charmingruby/winston/internal/agent/analyst"
	"github.com/charmingruby/winston/internal/agent/reviewer"
	"github.com/charmingruby/winston/internal/agent/strategist"
	"github.com/charmingruby/winston/internal/agent/winston"
	"github.com/charmingruby/winston/internal/tool"
	"github.com/firebase/genkit/go/ai/exp"
	"github.com/firebase/genkit/go/ai/exp/localstore"
)

func NewWinston(ctx context.Context, cfg *config.Config) (*winston.Winston, error) {
	google := provider.NewGoogle(cfg.GeminiAPIKey)
	modelName := cfg.GeminiModel
	ak := aikit.New(ctx, google)

	an, err := analyst.Define(ak, google, modelName)
	if err != nil {
		return nil, fmt.Errorf("define analyst: %w", err)
	}

	st, err := strategist.Define(ak, google, modelName)
	if err != nil {
		return nil, fmt.Errorf("define strategist: %w", err)
	}

	rv, err := reviewer.Define(ak, google, modelName)
	if err != nil {
		return nil, fmt.Errorf("define reviewer: %w", err)
	}

	w, err := winston.Define(
		ak,
		google,
		modelName,
		[]exp.AgentRef{an.Ref(), st.Ref(), rv.Ref()},
		[]aikit.Tool{tool.GetSquadRoaster()},
		exp.WithSessionStore(localstore.NewInMemorySessionStore[struct{}]()),
	)
	if err != nil {
		return nil, fmt.Errorf("define winston: %w", err)
	}

	return w, nil
}
