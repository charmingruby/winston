package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmingruby/stan/internal/platform/aikit"
	"github.com/charmingruby/stan/internal/platform/aikit/provider"
	"github.com/charmingruby/stan/internal/platform/config"
	"github.com/charmingruby/stan/internal/platform/o11y/logging"
	"github.com/firebase/genkit/go/ai"
	aix "github.com/firebase/genkit/go/ai/exp"
)

func main() {
	ctx := context.Background()

	log := logging.InitLogger()
	logging.WithLogger(ctx, log)

	if err := run(ctx); err != nil {
		log.ErrorContext(ctx, "winstom example failed", "error", err)
	}
}

func run(ctx context.Context) error {
	log := logging.LoggerFromContext(ctx)

	goal := strings.TrimSpace(strings.Join(os.Args[1:], " "))
	if goal == "" {
		return errors.New("goal cannot be blank; usage: go run ./examples/winstom <question>")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	google := provider.NewGoogle(cfg.GeminiAPIKey)
	ak := aikit.New(ctx, google)

	model := google.Model("gemini-3.5-flash-lite")

	agent, err := aikit.DefineAgent(
		ak,
		"winstom",
		model,
		nil,
		aix.InlinePrompt{
			ai.WithSystem(
				`You are Winstom, a friendly Overwatch assistant.
				Help with heroes, abilities, team compositions, and strategies.
				Keep answers short and clear.
				If you are unsure, say so.`,
			),
		},
	)
	if err != nil {
		return fmt.Errorf("define agent: %w", err)
	}

	response, err := agent.Run(ctx, goal)
	if err != nil {
		return fmt.Errorf("run agent: %w", err)
	}

	log.InfoContext(ctx,
		"agent ran successfully",
		"response", response,
	)

	return nil
}
