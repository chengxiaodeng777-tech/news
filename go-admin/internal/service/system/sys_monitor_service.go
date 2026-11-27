package system

import (
	"context"
	"database/sql"
	"os"
	"runtime"
	"time"

	"base/internal/config"
	"base/internal/db"
	systemmodel "base/internal/model/system"
	rediscache "base/internal/redis"
)

type MonitorService struct {
	dbResource  *db.Resource
	redisClient *rediscache.Client
	cfg         config.Config
	startTime   time.Time
}

func NewMonitor(resource *db.Resource, redisClient *rediscache.Client, cfg config.Config, startTime time.Time) *MonitorService {
	return &MonitorService{
		dbResource:  resource,
		redisClient: redisClient,
		cfg:         cfg,
		startTime:   startTime,
	}
}

func (s *MonitorService) GetMonitor(ctx context.Context) (*systemmodel.MonitorResponse, error) {
	dbStats := db.Stats(s.dbResource)
	dbStatus := "正常"
	if err := db.Ping(ctx, s.dbResource); err != nil {
		dbStatus = "异常"
	}

	redisStatus := "正常"
	redisPing, err := rediscache.Ping(ctx, s.redisClient)
	if err != nil {
		redisStatus = "异常"
		redisPing = 0
	}
	redisStats := rediscache.PoolStats(s.redisClient)

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	hostName, _ := os.Hostname()
	hostInfo := readHostInfo()

	lastGC := ""
	if mem.LastGC > 0 {
		lastGC = time.UnixMilli(int64(mem.LastGC / uint64(time.Millisecond))).Format("2006-01-02 15:04:05.000")
	}

	return &systemmodel.MonitorResponse{
		Overview: systemmodel.MonitorOverview{
			AppName:          "base",
			Environment:      s.cfg.App.Env,
			Address:          s.cfg.App.Addr,
			StartTime:        s.startTime.Format("2006-01-02 15:04:05.000"),
			Uptime:           time.Since(s.startTime).String(),
			RequestIDEnabled: true,
			CurrentTime:      time.Now().Format("2006-01-02 15:04:05.000"),
		},
		Database: systemmodel.MonitorDatabase{
			Status:             dbStatus,
			Driver:             dbStats.DriverName,
			OpenConnections:    dbStats.OpenConnections,
			InUseConnections:   dbStats.InUseConnections,
			IdleConnections:    dbStats.IdleConnections,
			WaitCount:          dbStats.WaitCount,
			WaitDuration:       dbStats.WaitDuration.String(),
			MaxOpenConnections: dbStats.MaxOpenConnections,
		},
		Redis: systemmodel.MonitorRedis{
			Status:       redisStatus,
			Address:      s.cfg.Redis.Addr,
			DB:           s.cfg.Redis.DB,
			PoolSize:     s.cfg.Redis.PoolSize,
			PoolHits:     redisStats.PoolHits,
			PoolMisses:   redisStats.PoolMisses,
			Timeouts:     redisStats.Timeouts,
			TotalConns:   redisStats.TotalConns,
			IdleConns:    redisStats.IdleConns,
			StaleConns:   redisStats.StaleConns,
			PingDuration: redisPing.String(),
		},
		Runtime: systemmodel.MonitorRuntime{
			GoVersion:     runtime.Version(),
			Goroutines:    runtime.NumGoroutine(),
			CPUs:          runtime.NumCPU(),
			MemoryAllocMB: mem.Alloc / 1024 / 1024,
			MemorySysMB:   mem.Sys / 1024 / 1024,
			HeapAllocMB:   mem.HeapAlloc / 1024 / 1024,
			HeapSysMB:     mem.HeapSys / 1024 / 1024,
			LastGC:        lastGC,
			NumGC:         mem.NumGC,
		},
		Host: systemmodel.MonitorHost{
			Hostname:      hostName,
			OS:            runtime.GOOS,
			Arch:          runtime.GOARCH,
			CPULoad:       hostInfo.CPULoad,
			TotalMemoryMB: hostInfo.TotalMemoryMB,
			FreeMemoryMB:  hostInfo.FreeMemoryMB,
			TotalDiskGB:   hostInfo.TotalDiskGB,
			FreeDiskGB:    hostInfo.FreeDiskGB,
		},
	}, nil
}

func sqlDBFromResource(resource *db.Resource) *sql.DB {
	return resource.SQLDB
}

type hostMetrics struct {
	CPULoad       float64
	TotalMemoryMB uint64
	FreeMemoryMB  uint64
	TotalDiskGB   uint64
	FreeDiskGB    uint64
}
