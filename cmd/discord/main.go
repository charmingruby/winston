package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmingruby/erebor/env"
	"github.com/charmingruby/erebor/o11y/logging"
	"github.com/charmingruby/winston/config"
	"github.com/charmingruby/winston/internal/discord"
	"github.com/charmingruby/winston/internal/workflow"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := logging.InitLogger()
	ctx = logging.WithLogger(ctx, log)

	if err := run(ctx); err != nil {
		log.ErrorContext(ctx, "winston discord bot failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := env.Load[config.Config]()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	w, err := workflow.NewWinston(ctx, cfg)
	if err != nil {
		return err
	}

	bot, err := discord.New(cfg.DiscordBotToken, func(ctx context.Context, question string) (string, error) {
		return w.Run(ctx, question)
	})
	if err != nil {
		return err
	}

	return bot.ListenToMentions(ctx)
}
