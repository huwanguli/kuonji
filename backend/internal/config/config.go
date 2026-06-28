package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Expire time.Duration `mapstructure:"expire"`
}

type LogConfig struct {
	Level   string `mapstructure:"level"`
	File    string `mapstructure:"file"`
	MaxAge  int    `mapstructure:"max_age"`
	MaxSize int    `mapstructure:"max_size"`
}

type UploadConfig struct {
	Path        string   `mapstructure:"path"`
	MaxSize     int      `mapstructure:"max_size"`
	AllowedExts []string `mapstructure:"allowed_exts"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	Upload   UploadConfig   `mapstructure:"upload"`
}

var Cfg *Config

func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	v.AutomaticEnv()
	v.SetEnvPrefix("ZBLOG")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	bindEnvs(v, "server.port", "server.mode")
	bindEnvs(v, "database.host", "database.port", "database.user", "database.password", "database.dbname")
	bindEnvs(v, "jwt.secret")
	bindEnvs(v, "log.level")

	var cfg Config
	decoderOpt := viper.DecodeHook(mapstructure.StringToTimeDurationHookFunc())
	if err := v.Unmarshal(&cfg, decoderOpt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	Cfg = &cfg
	return &cfg, nil
}

func bindEnvs(v *viper.Viper, keys ...string) {
	for _, k := range keys {
		v.BindEnv(k)
	}
}
