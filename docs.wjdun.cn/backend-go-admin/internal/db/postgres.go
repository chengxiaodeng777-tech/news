package db

import (
	"context"
	"database/sql"
	"net"
	"net/url"
	"strconv"
	"time"

	"base/ent"
	"base/internal/config"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type Resource struct {
	Client *ent.Client
	SQLDB  *sql.DB
}

type DBStats struct {
	DriverName         string
	OpenConnections    int
	InUseConnections   int
	IdleConnections    int
	WaitCount          int64
	WaitDuration       time.Duration
	MaxOpenConnections int
}

func Open(ctx context.Context, cfg config.DatabaseConfig) (*Resource, error) {
	connConfig, err := pgx.ParseConfig(buildDatabaseURL(cfg))
	if err != nil {
		return nil, err
	}

	sqlDB := stdlib.OpenDB(*connConfig)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeMinutes) * time.Minute)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	driver := entsql.OpenDB(dialect.Postgres, sqlDB)
	return &Resource{
		Client: ent.NewClient(ent.Driver(driver)),
		SQLDB:  sqlDB,
	}, nil
}

func Migrate(ctx context.Context, resource *Resource) error {
	if err := resource.Client.Schema.Create(ctx); err != nil {
		return err
	}
	return migrateArticleSort(ctx, resource.SQLDB)
}

func Stats(resource *Resource) DBStats {
	stats := resource.SQLDB.Stats()
	return DBStats{
		DriverName:         "postgres",
		OpenConnections:    stats.OpenConnections,
		InUseConnections:   stats.InUse,
		IdleConnections:    stats.Idle,
		WaitCount:          stats.WaitCount,
		WaitDuration:       stats.WaitDuration,
		MaxOpenConnections: stats.MaxOpenConnections,
	}
}

func Ping(ctx context.Context, resource *Resource) error {
	return resource.SQLDB.PingContext(ctx)
}

func Close(resource *Resource) error {
	if resource == nil {
		return nil
	}
	return resource.Client.Close()
}

func buildDatabaseURL(cfg config.DatabaseConfig) string {
	query := url.Values{}
	query.Set("sslmode", cfg.SSLMode)
	if cfg.TimeZone != "" {
		query.Set("timezone", cfg.TimeZone)
	}

	user := url.User(cfg.User)
	if cfg.Password != "" {
		user = url.UserPassword(cfg.User, cfg.Password)
	}

	return (&url.URL{
		Scheme:   "postgres",
		User:     user,
		Host:     net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Path:     cfg.Name,
		RawQuery: query.Encode(),
	}).String()
}
