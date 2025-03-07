package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	EnvironmentVariables EnvironmentVariables `yaml:"EnvironmentVariables"`
	Server               Server               `yaml:"Server"`
	Logs                 Logs                 `yaml:"Logs"`
	Cors                 Cors                 `yaml:"Cors"`
}

type EnvironmentVariables struct {
	Environment string `yaml:"Environment"`
}

type Server struct {
	Addr    string        `yaml:"Addr"`
	Port    string        `yaml:"Port"`
	Timeout time.Duration `yaml:"Timeout"`
}

type Logs struct {
	Path       string `yaml:"Path"`
	Level      string `yaml:"Level"`
	MaxAge     int    `yaml:"MaxAge"`
	MaxBackups int    `yaml:"MaxBackups"`
}

type Cors struct {
	AllowedOrigins []string `yaml:"AllowedOrigins"`
}

func ReadConfig(cfgName, cfgType, cfgPath string) (*Config, error) {
	var cfg Config

	viper.SetConfigName(cfgName)
	viper.SetConfigType(cfgType)
	viper.AddConfigPath(cfgPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
