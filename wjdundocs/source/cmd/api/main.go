package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"base/internal/bootstrap"
	internalcasbin "base/internal/casbin"
	"base/internal/config"
	"base/internal/db"
	"base/internal/logger"
	"base/internal/redis"
	"base/internal/router"
)

func main() {
	startTime := time.Now()
	appCtx, cancelApp := context.WithCancel(context.WithValue(context.Background(), "app_start_time", startTime))
	defer cancelApp()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log, err := logger.New(cfg.Log)
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()

	database, err := db.Open(appCtx, cfg.Database)
	if err != nil {
		log.Fatal("打开数据库失败", zapError(err))
	}
	defer db.Close(database)
	if err := db.Migrate(appCtx, database); err != nil {
		log.Fatal("执行数据库迁移失败", zapError(err))
	}
	if cfg.App.AutoSeed {
		if err := bootstrap.InitSystemData(appCtx, database.Client); err != nil {
			log.Fatal("初始化系统基础数据失败", zapError(err))
		}
	}

	enforcer, err := internalcasbin.NewEnforcer(database.Client)
	if err != nil {
		log.Fatal("初始化 Casbin 失败", zapError(err))
	}

	redisClient, err := redis.Open(appCtx, cfg.Redis)
	if err != nil {
		log.Fatal("打开 Redis 失败", zapError(err))
	}
	defer redisClient.Close()

	r := router.New(appCtx, cfg, log, router.Dependencies{
		Database: database,
		Redis:    redisClient,
		Enforcer: enforcer,
	})
	srv := &http.Server{
		Addr:    cfg.App.Addr,
		Handler: r,
		BaseContext: func(_ net.Listener) context.Context {
			return appCtx
		},
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Info("服务启动成功", logger.String("addr", cfg.App.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-signalCtx.Done():
		log.Info("收到退出信号，开始优雅关闭")
	case err := <-serverErr:
		if err != nil {
			log.Fatal("服务运行失败", zapError(err))
		}
		return
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(
		context.Background(),
		time.Duration(cfg.App.ShutdownTimeoutSeconds)*time.Second,
	)
	defer cancelShutdown()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("服务关闭失败", zapError(err))
		return
	}

	cancelApp()
	if err := <-serverErr; err != nil {
		log.Error("服务关闭后监听异常", zapError(err))
	}
	log.Info("服务已优雅关闭")
}

func zapError(err error) logger.Field {
	return logger.Error(err)
}
