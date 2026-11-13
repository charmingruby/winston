package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmingruby/winstom/internal/agent/analyst"
	"github.com/charmingruby/winstom/internal/agent/reviewer"
	"github.com/charmingruby/winstom/internal/agent/strategist"
	"github.com/charmingruby/winstom/internal/agent/winstom"
	"github.com/charmingruby/winstom/internal/platform/aikit"
	"github.com/charmingruby/winstom/internal/platform/aikit/provider"
	"github.com/charmingruby/winstom/internal/platform/config"
	"github.com/charmingruby/winstom/internal/platform/httpx"
	"github.com/charmingruby/winstom/internal/platform/o11y/logging"
	aix "github.com/firebase/genkit/go/ai/exp"
	"github.com/firebase/genkit/go/ai/exp/localstore"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := logging.InitLogger()
	ctx = logging.WithLogger(ctx, log)

	if err := run(ctx); err != nil {
		log.ErrorContext(ctx, "winstom example failed", "error", err)
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
	modelName := "gemini-3.5-flash-lite"

	an, err := analyst.Define(ak, google, modelName)
	if err != nil {
		return fmt.Errorf("define analyst: %w", err)
	}

	st, err := strategist.Define(ak, google, modelName)
	if err != nil {
		return fmt.Errorf("define strategist: %w", err)
	}

	rv, err := reviewer.Define(ak, google, modelName)
	if err != nil {
		return fmt.Errorf("define reviewer: %w", err)
	}

	w, err := winstom.Define(
		ak,
		google,
		modelName,
		[]aix.AgentRef{
			an.Ref(),
			st.Ref(),
			rv.Ref(),
		},
		aix.WithSessionStore(localstore.NewInMemorySessionStore[struct{}]()),
	)
	if err != nil {
		return fmt.Errorf("define winstom: %w", err)
	}

	srv, r := httpx.NewServer(cfg.Port)
	srv.WriteTimeout = 4 * time.Minute

	r.Post("/agents/winstom", func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Minute)
		defer cancel()

		w.HandleHTTP(rw, r.WithContext(ctx))
	})

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.Start()
	}()

	log.InfoContext(ctx, "winstom server listening",
		"addr", srv.Addr,
	)

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return srv.Close(ctx)
	}
}
