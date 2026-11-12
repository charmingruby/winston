package main

import (
	"context"
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

	message := strings.TrimSpace(strings.Join(os.Args[1:], " "))
	if message == "" {
		log.Error(
			"message cannot be blank",
			"usage", "go run examples/basicagent/main.go <question>",
		)

		return
	}

	cfg, err := config.Load()
	if err != nil {
		log.Error(
			"error loading config",
			"error", err,
		)

		return
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
		log.Error(
			"error defining agent",
			"error", err,
		)

		return
	}

	response, err := agent.Run(ctx, message)
	if err != nil {
		log.Error(
			"error running agent",
			"error", err,
		)

		return
	}

	log.Info(
		"agent ran successfully",
		"response", response,
	)
}
