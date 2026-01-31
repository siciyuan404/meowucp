//go:build tools
// +build tools

package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

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

func loadConfig(path string) (*DatabaseConfig, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg struct {
		Database DatabaseConfig
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg.Database, nil
}

func getDSN(cfg *DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
}

func main() {
	cfg, err := loadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dsn := getDSN(cfg)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connected successfully")

	migrations := []string{
		"migrations/000_create_extension.sql",
		"migrations/001_init.sql",
		"migrations/002_ucp.sql",
		"migrations/003_ucp_webhook_audit.sql",
		"migrations/004_ucp_webhook_replay_queue.sql",
		"migrations/005_ucp_webhook_alerts.sql",
	}

	for _, path := range migrations {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", path, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatalf("Failed to execute migration %s: %v", path, err)
		}
	}

	log.Println("Migrations executed successfully")
}
