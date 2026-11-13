package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmingruby/stan/internal/platform/aikit"
	"github.com/charmingruby/stan/internal/platform/aikit/provider"
	"github.com/charmingruby/stan/internal/platform/config"
	"github.com/charmingruby/stan/internal/platform/httpx"
	"github.com/charmingruby/stan/internal/platform/o11y/logging"
	"github.com/firebase/genkit/go/ai"
	aix "github.com/firebase/genkit/go/ai/exp"
	"github.com/firebase/genkit/go/ai/exp/localstore"
	"github.com/firebase/genkit/go/genkit"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := logging.InitLogger()
	ctx = logging.WithLogger(ctx, log)

	if err := run(ctx); err != nil {
		log.ErrorContext(ctx, "agent server failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	log := logging.LoggerFromContext(ctx)

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	google := provider.NewGoogle(cfg.GeminiAPIKey)
	ak := aikit.New(ctx, google)
	model := google.Model("gemini-3.5-flash-lite")

	agent, err := aikit.DefineAgent(ak, "assistant",
		model,
		nil,
		aix.InlinePrompt{
			ai.WithSystem(`You are Winstom, a friendly Overwatch assistant.
Help with heroes, abilities, team compositions, counters, maps, modes, and gameplay strategies.
Adapt recommendations to the supplied rank, role, player preferences, and match context.
For specific advice, ask focused questions when essential context is missing; do not invent it.
For general questions, give useful guidance without demanding unnecessary details.
You have no tools, live match data, profile access, replay access, or current patch verification.
Do not invent abilities, statistics, patch changes, or current-meta rankings. State uncertainty and distinguish general principles from patch-dependent advice.
Offer practical priorities, hero alternatives, and trade-offs. Encourage player comfort, constructive communication, and teamwork; never blame teammates or promise wins.
Keep answers short and clear, in Portuguese.`),
		},
		aix.WithSessionStore(localstore.NewInMemorySessionStore[struct{}]()),
	)
	if err != nil {
		return fmt.Errorf("define agent: %w", err)
	}

	srv, router := httpx.NewServer(cfg.Port)
	srv.WriteTimeout = 2 * time.Minute

	agentHandler := genkit.Handler(agent.Genkit)

	router.Post("/agents/assistant", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
		defer cancel()

		agentHandler(w, r.WithContext(ctx))
	})

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.Start()
	}()

	log.InfoContext(ctx, "agent server listening",
		"addr", srv.Addr,
	)

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Close(shutdownCtx)
	}
}
