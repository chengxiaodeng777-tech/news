package redis

import (
	"context"
	"time"

	"base/internal/config"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	*goredis.Client
}

type Stats struct {
	PoolHits   uint32
	PoolMisses uint32
	Timeouts   uint32
	TotalConns uint32
	IdleConns  uint32
	StaleConns uint32
}

func Open(ctx context.Context, cfg config.RedisConfig) (*Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  time.Duration(cfg.DialTimeoutSeconds) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
		PoolSize:     cfg.PoolSize,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}

	return &Client{Client: client}, nil
}

func Ping(ctx context.Context, client *Client) (time.Duration, error) {
	start := time.Now()
	if err := client.Client.Ping(ctx).Err(); err != nil {
		return 0, err
	}
	return time.Since(start), nil
}

func PoolStats(client *Client) Stats {
	stats := client.Client.PoolStats()
	return Stats{
		PoolHits:   stats.Hits,
		PoolMisses: stats.Misses,
		Timeouts:   stats.Timeouts,
		TotalConns: stats.TotalConns,
		IdleConns:  stats.IdleConns,
		StaleConns: stats.StaleConns,
	}
}
