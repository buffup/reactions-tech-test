package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/buffup/reactions-tech-test/api"
	"github.com/buffup/reactions-tech-test/daemon"
	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	RedisHost string `envconfig:"REDIS_HOST" required:"true"`
	RedisPort int    `envconfig:"REDIS_PORT" required:"true"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	redisClient := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)})

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		srv := &http.Server{
			Addr: ":8080",
			Handler: &api.API{
				Cache:              redisClient,
				AvailableReactions: []string{"like", "love", "haha", "wow", "sad", "angry"},
			},
		}

		go func() {
			<-ctx.Done()
			srv.Shutdown(context.Background())
		}()

		return srv.ListenAndServe()
	})

	eg.Go(func() error {
		daemon := &daemon.Daemon{
			Cache:        redisClient,
			Pubsub:       &PubSub{},
			SendInterval: time.Second,
		}

		return daemon.Run(ctx)
	})

	if err := eg.Wait(); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
