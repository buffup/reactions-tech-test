package pubsub

import (
	"context"
	"log/slog"
)

type PubSub struct{}

func (p *PubSub) Publish(ctx context.Context, channel string, payload any) error {
	slog.Info("New message", "channel", channel, "payload", payload)
	return nil
}
