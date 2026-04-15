package config

import "github.com/spf13/viper"

var AppConfig *Config

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Postgres PostgresConfig `mapstructure:"postgres"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

func LoadConfig(path *string) (*Config, error) {
	v := viper.New()

	if path != nil {
		v.SetConfigFile(*path)
	} else {
		v.SetConfigFile("app.yaml")
	}
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if cfg.Server.Port != "" && cfg.Server.Port[0] != ':' {
		cfg.Server.Port = ":" + cfg.Server.Port
	}

	AppConfig = &cfg
	return AppConfig, nil
}
