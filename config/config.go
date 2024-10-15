package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AuthenticationConfiguration AuthenticationConfiguration `yaml:"AuthenticationConfiguration"`
	EnvironmentVariables        EnvironmentVariables        `yaml:"EnvironmentVariables"`
	Server                      Server                      `yaml:"Server"`
	Logs                        Logs                        `yaml:"Logs"`
}

type AuthenticationConfiguration struct {
	AccessSecretKey string `yaml:"AccessSecretKey"`
}

type EnvironmentVariables struct {
	Environment string `yaml:"Environment"`
}

type Server struct {
	Type string `yaml:"Type"`
	Port string `yaml:"Port"`
}

type Logs struct {
	Path       string `yaml:"Path"`
	Level      string `yaml:"Level"`
	MaxAge     int    `yaml:"MaxAge"`
	MaxBackups int    `yaml:"MaxBackups"`
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
