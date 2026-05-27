package config

import (
    "fmt"
    "time"

    "github.com/fsnotify/fsnotify"
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
    Port         int           `mapstructure:"port"`
    Mode         string        `mapstructure:"mode"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
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
    Secret      string `mapstructure:"secret"`
    Issuer      string `mapstructure:"issuer"`
    ExpireHours int    `mapstructure:"expire_hours"`
}

type LogConfig struct {
    Level      string `mapstructure:"level"`
    Filename   string `mapstructure:"filename"`
    MaxSize    int    `mapstructure:"max_size"`
    MaxBackups int    `mapstructure:"max_backups"`
    MaxAge     int    `mapstructure:"max_age"`
    Compress   bool   `mapstructure:"compress"`
}

var globalConfig *Config

func LoadConfig(configPath string) (*Config, error) {
    v := viper.New()
    v.SetConfigFile(configPath)
    v.SetConfigType("yaml")

    if err := v.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    v.WatchConfig()
    v.OnConfigChange(func(e fsnotify.Event) {
        var newConfig Config
        if err := v.Unmarshal(&newConfig); err != nil {
            return
        }
        globalConfig = &newConfig
    })

    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }

    globalConfig = &cfg
    return globalConfig, nil
}

func GetConfig() *Config {
    return globalConfig
}