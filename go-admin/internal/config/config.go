package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	configDir      = "configs"
	baseConfigFile = "config.yaml"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	Log      LogConfig      `yaml:"log"`
}

type AppConfig struct {
	Env                    string `yaml:"env"`
	Addr                   string `yaml:"addr"`
	ShutdownTimeoutSeconds int    `yaml:"shutdown_timeout_seconds"`
	AutoSeed               bool   `yaml:"auto_seed"`
}

type DatabaseConfig struct {
	Host                   string `yaml:"host"`
	Port                   int    `yaml:"port"`
	User                   string `yaml:"user"`
	Password               string `yaml:"password"`
	Name                   string `yaml:"name"`
	SSLMode                string `yaml:"ssl_mode"`
	TimeZone               string `yaml:"time_zone"`
	MaxOpenConns           int    `yaml:"max_open_conns"`
	MaxIdleConns           int    `yaml:"max_idle_conns"`
	ConnMaxLifetimeMinutes int    `yaml:"conn_max_lifetime_minutes"`
}

type RedisConfig struct {
	Addr                string `yaml:"addr"`
	Password            string `yaml:"password"`
	DB                  int    `yaml:"db"`
	DialTimeoutSeconds  int    `yaml:"dial_timeout_seconds"`
	ReadTimeoutSeconds  int    `yaml:"read_timeout_seconds"`
	WriteTimeoutSeconds int    `yaml:"write_timeout_seconds"`
	PoolSize            int    `yaml:"pool_size"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

type LogConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	Directory  string `yaml:"directory"`
	Filename   string `yaml:"filename"`
	MaxSizeMB  int    `yaml:"max_size_mb"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAgeDays int    `yaml:"max_age_days"`
	Compress   bool   `yaml:"compress"`
}

func Load() (Config, error) {
	cfg := defaultConfig()

	if err := mergeFile(&cfg, baseConfigFile); err != nil {
		return Config{}, err
	}

	env := getEnv("APP_ENV", cfg.App.Env)
	if env == "" {
		env = "dev"
	}

	cfg.App.Env = env
	if err := mergeFile(&cfg, envConfigFile(env)); err != nil {
		return Config{}, err
	}

	cfg.App.Addr = getEnv("APP_ADDR", cfg.App.Addr)

	if cfg.App.Addr == "" {
		return Config{}, errors.New("配置项 app.addr 不能为空")
	}
	if cfg.App.ShutdownTimeoutSeconds <= 0 {
		return Config{}, errors.New("配置项 app.shutdown_timeout_seconds 必须大于 0")
	}
	if cfg.Database.Host == "" {
		return Config{}, errors.New("配置项 database.host 不能为空")
	}
	if cfg.Database.Port <= 0 {
		return Config{}, errors.New("配置项 database.port 必须大于 0")
	}
	if cfg.Database.User == "" {
		return Config{}, errors.New("配置项 database.user 不能为空")
	}
	if cfg.Database.Name == "" {
		return Config{}, errors.New("配置项 database.name 不能为空")
	}
	if cfg.Database.SSLMode == "" {
		return Config{}, errors.New("配置项 database.ssl_mode 不能为空")
	}
	if cfg.Database.MaxOpenConns <= 0 {
		return Config{}, errors.New("配置项 database.max_open_conns 必须大于 0")
	}
	if cfg.Database.MaxIdleConns <= 0 {
		return Config{}, errors.New("配置项 database.max_idle_conns 必须大于 0")
	}
	if cfg.Database.ConnMaxLifetimeMinutes <= 0 {
		return Config{}, errors.New("配置项 database.conn_max_lifetime_minutes 必须大于 0")
	}
	if cfg.Redis.Addr == "" {
		return Config{}, errors.New("配置项 redis.addr 不能为空")
	}
	if cfg.Redis.DialTimeoutSeconds <= 0 {
		return Config{}, errors.New("配置项 redis.dial_timeout_seconds 必须大于 0")
	}
	if cfg.Redis.ReadTimeoutSeconds <= 0 {
		return Config{}, errors.New("配置项 redis.read_timeout_seconds 必须大于 0")
	}
	if cfg.Redis.WriteTimeoutSeconds <= 0 {
		return Config{}, errors.New("配置项 redis.write_timeout_seconds 必须大于 0")
	}
	if cfg.Redis.PoolSize <= 0 {
		return Config{}, errors.New("配置项 redis.pool_size 必须大于 0")
	}
	if cfg.JWT.Secret == "" {
		return Config{}, errors.New("配置项 jwt.secret 不能为空")
	}
	if cfg.Log.Format == "" {
		return Config{}, errors.New("配置项 log.format 不能为空")
	}
	if cfg.Log.Output == "" {
		return Config{}, errors.New("配置项 log.output 不能为空")
	}
	if !isValidLogOutput(cfg.Log.Output) {
		return Config{}, fmt.Errorf("配置项 log.output 不合法，可选值为 console、file、both")
	}
	if !isValidLogFormat(cfg.Log.Format) {
		return Config{}, fmt.Errorf("配置项 log.format 不合法，可选值为 console、json")
	}
	if cfg.Log.Directory == "" {
		return Config{}, errors.New("配置项 log.directory 不能为空")
	}
	if cfg.Log.Filename == "" {
		return Config{}, errors.New("配置项 log.filename 不能为空")
	}
	if cfg.Log.MaxSizeMB <= 0 {
		return Config{}, errors.New("配置项 log.max_size_mb 必须大于 0")
	}
	if cfg.Log.MaxAgeDays <= 0 {
		return Config{}, errors.New("配置项 log.max_age_days 必须大于 0")
	}

	return cfg, nil
}

func defaultConfig() Config {
	return Config{
		App: AppConfig{
			Env:                    "dev",
			Addr:                   ":8080",
			ShutdownTimeoutSeconds: 30,
			AutoSeed:               true,
		},
		Database: DatabaseConfig{
			Host:                   "localhost",
			Port:                   5432,
			User:                   "postgres",
			Password:               "postgres",
			Name:                   "base",
			SSLMode:                "disable",
			TimeZone:               "Asia/Shanghai",
			MaxOpenConns:           10,
			MaxIdleConns:           10,
			ConnMaxLifetimeMinutes: 30,
		},
		Redis: RedisConfig{
			Addr:                "localhost:6379",
			DB:                  0,
			DialTimeoutSeconds:  5,
			ReadTimeoutSeconds:  3,
			WriteTimeoutSeconds: 3,
			PoolSize:            10,
		},
		Log: LogConfig{
			Level:      "info",
			Format:     "console",
			Output:     "both",
			Directory:  "logs",
			Filename:   "app.log",
			MaxSizeMB:  10,
			MaxBackups: 10,
			MaxAgeDays: 7,
			Compress:   true,
		},
	}
}

func envConfigFile(env string) string {
	return "config-" + env + ".yaml"
}

func mergeFile(cfg *Config, path string) error {
	data, err := os.ReadFile(filepath.Join(configDir, path))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, cfg)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func isValidLogOutput(output string) bool {
	switch strings.ToLower(output) {
	case "console", "file", "both":
		return true
	default:
		return false
	}
}

func isValidLogFormat(format string) bool {
	switch strings.ToLower(format) {
	case "console", "json":
		return true
	default:
		return false
	}
}
