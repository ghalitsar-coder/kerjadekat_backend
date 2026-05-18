package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds runtime settings loaded from environment and optional .env file.
type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	Env        string `mapstructure:"ENV"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBSslMode  string `mapstructure:"DB_SSLMODE"`

	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	RabbitMQURL string `mapstructure:"RABBITMQ_URL"`

	WorkerLocationTTLMinutes int `mapstructure:"WORKER_LOCATION_TTL_MINUTES"`
	MatchRadiusMeters        int `mapstructure:"MATCH_RADIUS_METERS"`
	MatchTimerSeconds        int `mapstructure:"MATCH_TIMER_SECONDS"`
	MatchMaxRounds           int `mapstructure:"MATCH_MAX_ROUNDS"`

	JWTSecret              string        `mapstructure:"JWT_SECRET"`
	JWTAccessExpiryMinutes int           `mapstructure:"JWT_ACCESS_EXPIRY_MINUTES"`
	JWTRefreshExpiryDays   int           `mapstructure:"JWT_REFRESH_EXPIRY_DAYS"`
	JWTAccessExpiry        time.Duration `mapstructure:"-"`
	JWTRefreshExpiry       time.Duration `mapstructure:"-"`

	XenditAPIKey          string `mapstructure:"XENDIT_API_KEY"`
	XenditCallbackToken   string `mapstructure:"XENDIT_CALLBACK_TOKEN"`
	OTPLogFile            string `mapstructure:"OTP_LOG_FILE"`
}

// Load reads configuration from OS environment variables and, if present, a .env file.
// Pass empty path to look for ".env" in the current working directory.
func Load(envFileDir string) (*Config, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("ENV", "development")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5435")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", "6379")
	v.SetDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	v.SetDefault("WORKER_LOCATION_TTL_MINUTES", 5)
	v.SetDefault("MATCH_RADIUS_METERS", 5000)
	v.SetDefault("MATCH_TIMER_SECONDS", 60)
	v.SetDefault("MATCH_MAX_ROUNDS", 5)
	v.SetDefault("JWT_ACCESS_EXPIRY_MINUTES", 15)
	v.SetDefault("JWT_REFRESH_EXPIRY_DAYS", 7)
	v.SetDefault("OTP_LOG_FILE", "tmp/otp.log")

	// Unmarshal only sees keys present in Viper; bind env vars explicitly so
	// OS environment works even when no .env file is used.
	envKeys := []string{
		"SERVER_PORT",
		"ENV",
		"DB_HOST",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_PORT",
		"DB_SSLMODE",
		"REDIS_HOST",
		"REDIS_PORT",
		"REDIS_PASSWORD",
		"RABBITMQ_URL",
		"WORKER_LOCATION_TTL_MINUTES",
		"MATCH_RADIUS_METERS",
		"MATCH_TIMER_SECONDS",
		"MATCH_MAX_ROUNDS",
		"JWT_SECRET",
		"JWT_ACCESS_EXPIRY_MINUTES",
		"JWT_REFRESH_EXPIRY_DAYS",
		"XENDIT_API_KEY",
		"XENDIT_CALLBACK_TOKEN",
		"OTP_LOG_FILE",
	}
	for _, k := range envKeys {
		if err := v.BindEnv(k); err != nil {
			return nil, fmt.Errorf("bind env %q: %w", k, err)
		}
	}

	envPath := ".env"
	if envFileDir != "" {
		envPath = filepath.Clean(filepath.Join(envFileDir, ".env"))
	}
	v.SetConfigType("env")
	if fi, err := os.Stat(envPath); err == nil && !fi.IsDir() {
		v.SetConfigFile(envPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config file %q: %w", envPath, err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	cfg.JWTAccessExpiry = time.Duration(cfg.JWTAccessExpiryMinutes) * time.Minute
	cfg.JWTRefreshExpiry = time.Duration(cfg.JWTRefreshExpiryDays) * 24 * time.Hour

	return &cfg, nil
}
