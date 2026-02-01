package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Log      LogConfig
	Cache    CacheConfig
	Queue    QueueConfig
	UCP      UCPConfig
}

type ServerConfig struct {
	Port         int
	Mode         string
	MetricsToken string `mapstructure:"metrics_token"`
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

type JWTConfig struct {
	Secret     string
	ExpireHour int
}

type LogConfig struct {
	Level  string
	Format string
	Output string
}

type CacheConfig struct {
	ExpireSeconds int
}

type QueueConfig struct {
	StreamKey     string
	ConsumerGroup string
}

type UCPConfig struct {
	Links           []UCPLinkConfig  `mapstructure:"links"`
	ContinueURLBase string           `mapstructure:"continue_url_base"`
	Webhook         UCPWebhookConfig `mapstructure:"webhook"`
}

type UCPLinkConfig struct {
	Type  string `mapstructure:"type"`
	URL   string `mapstructure:"url"`
	Title string `mapstructure:"title"`
}

type UCPWebhookConfig struct {
	JWKSetURL           string `mapstructure:"jwk_set_url"`
	ClockSkewSeconds    int    `mapstructure:"clock_skew_seconds"`
	DeliveryURL         string `mapstructure:"delivery_url"`
	DeliveryTimeoutSec  int    `mapstructure:"delivery_timeout_sec"`
	SkipSignatureVerify bool   `mapstructure:"skip_signature_verify"`
	AlertMinAttempts    int    `mapstructure:"alert_min_attempts"`
	AlertDedupeSeconds  int    `mapstructure:"alert_dedupe_seconds"`
}

func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func GetDSN(cfg *DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
}

func GetRedisAddr(cfg *RedisConfig) string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}
