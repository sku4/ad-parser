package configs

import (
	"context"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Profiles  []string `mapstructure:"profiles"`
	Parser    `mapstructure:"parser"`
	Tarantool `mapstructure:"tarantool"`
}

type Parser struct {
	CheckTime           time.Duration `mapstructure:"check_time"`
	TooManyReqLimit     int           `mapstructure:"too_many_requests_limit"`
	DownloadWorkerCount int           `mapstructure:"download_worker_count"`
	CleanTime           time.Duration `mapstructure:"clean_time"`
}

type Tarantool struct {
	Servers           []string      `mapstructure:"servers"`
	User              string        `mapstructure:"user"`
	Password          string        `mapstructure:"password"`
	Timeout           time.Duration `mapstructure:"timeout"`
	ReconnectInterval time.Duration `mapstructure:"reconnect_interval"`
}

func Init() (*Config, error) {
	mainViper := viper.New()
	mainViper.AddConfigPath("configs")
	if err := mainViper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := mainViper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type configKey struct{}

func Set(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, configKey{}, cfg)
}

func Get(ctx context.Context) *Config {
	contextConfig, _ := ctx.Value(configKey{}).(*Config)

	return contextConfig
}
