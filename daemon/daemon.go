package daemon

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/buffup/reactions-tech-test/pubsub"
	"github.com/redis/go-redis/v9"
)

type Daemon struct {
	Cache        *redis.Client
	Pubsub       *pubsub.PubSub
	SendInterval time.Duration
}

func (d *Daemon) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-time.After(d.SendInterval):
			if err := d.sendReactions(ctx); err != nil {
				return err
			}
		}
	}
}

func (d *Daemon) sendReactions(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, d.SendInterval)
	defer cancel()

	keys, err := d.Cache.Keys(ctx, "livestreams:*:reactions:*").Result()
	if err != nil {
		return err
	}

	var snapshots = make(map[string]*ReactionSnapshot)
	for _, key := range keys {
		var (
			livestream = strings.Split(key, ":")[1]
			reaction   = strings.Split(key, ":")[3]
		)

		count, err := d.Cache.Get(ctx, key).Int()
		if err != nil {
			return err
		}

		if err := d.Cache.Del(ctx, key).Err(); err != nil {
			return err
		}

		if _, ok := snapshots[livestream]; !ok {
			snapshots[livestream] = &ReactionSnapshot{
				Livestream: livestream,
				Timestamp:  time.Now(),
				Reactions:  make(map[string]int),
			}
		}

		snapshots[livestream].Reactions[reaction] = count
	}

	for _, snapshot := range snapshots {
		if err := d.Pubsub.Publish(ctx, fmt.Sprintf("reactions.%s", snapshot.Livestream), snapshot); err != nil {
			return err
		}
	}

	return nil
}

type ReactionSnapshot struct {
	Livestream string         `json:"livestream"`
	Timestamp  time.Time      `json:"timestamp"`
	Reactions  map[string]int `json:"reactions"`
}
