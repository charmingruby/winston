package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmingruby/erebor/env"
	"github.com/charmingruby/erebor/httpx"
	"github.com/charmingruby/erebor/o11y/logging"
	"github.com/charmingruby/winston/config"
	"github.com/charmingruby/winston/internal/workflow"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := logging.InitLogger()
	ctx = logging.WithLogger(ctx, log)

	if err := run(ctx); err != nil {
		log.ErrorContext(ctx, "winston example failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	log := logging.LoggerFromContext(ctx)

	cfg, err := env.Load[config.Config]()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	w, err := workflow.NewWinston(ctx, cfg)
	if err != nil {
		return err
	}

	srv, r := httpx.NewServer(cfg.Port)
	srv.WriteTimeout = 4 * time.Minute

	r.Post("/agents/winston", func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Minute)
		defer cancel()

		w.HandleHTTP(rw, r.WithContext(ctx))
	})

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.Start()
	}()

	log.InfoContext(ctx, "winston server listening",
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
