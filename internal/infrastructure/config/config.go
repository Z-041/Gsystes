package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Log       LogConfig       `mapstructure:"log"`
	Upload    UploadConfig    `mapstructure:"upload"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	CORS      CORSConfig      `mapstructure:"cors"`
	Admin     AdminConfig     `mapstructure:"admin"`
}

type AdminConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Nickname string `mapstructure:"nickname"`
	Email    string `mapstructure:"email"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	MaxBodySize  int64         `mapstructure:"max_body_size"`
}

type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	Charset         string `mapstructure:"charset"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		d.Username, d.Password, d.Host, d.Port, d.Database, d.Charset)
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type JWTConfig struct {
	Secret        string `mapstructure:"secret"`
	Issuer        string `mapstructure:"issuer"`
	ExpireHours   int    `mapstructure:"expire_hours"`
	SigningMethod string `mapstructure:"signing_method"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type UploadConfig struct {
	Dir         string   `mapstructure:"dir"`
	MaxSize     int      `mapstructure:"max_size"`
	AllowedExts []string `mapstructure:"allowed_exts"`
}

type RateLimitConfig struct {
	Enabled      bool          `mapstructure:"enabled"`
	DefaultRate  int           `mapstructure:"default_rate"`
	DefaultBurst int           `mapstructure:"default_burst"`
	Window       time.Duration `mapstructure:"window"`
	LoginRate    int           `mapstructure:"login_rate"`
	LoginBurst   int           `mapstructure:"login_burst"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

var (
	globalConfig *Config
	configMu     sync.RWMutex
)

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	v.SetEnvPrefix("GSYSTES")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		var newConfig Config
		if err := v.Unmarshal(&newConfig); err != nil {
			return
		}
		configMu.Lock()
		globalConfig = &newConfig
		configMu.Unlock()
	})

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	applyEnvOverrides(&cfg)

	configMu.Lock()
	globalConfig = &cfg
	configMu.Unlock()

	return globalConfig, nil
}

func envOrDefault(envKey, fallback string) string {
	if val := os.Getenv(envKey); val != "" {
		return val
	}
	return fallback
}

func applyEnvOverrides(cfg *Config) {
	cfg.Database.Password = envOrDefault("GSYSTES_DATABASE_PASSWORD", cfg.Database.Password)
	cfg.JWT.Secret = envOrDefault("GSYSTES_JWT_SECRET", cfg.JWT.Secret)
	cfg.Redis.Password = envOrDefault("GSYSTES_REDIS_PASSWORD", cfg.Redis.Password)
	cfg.Admin.Password = envOrDefault("GSYSTES_ADMIN_PASSWORD", cfg.Admin.Password)
}

func GetConfig() *Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return globalConfig
}

func SetConfigForTesting(cfg *Config) {
	configMu.Lock()
	globalConfig = cfg
	configMu.Unlock()
}
