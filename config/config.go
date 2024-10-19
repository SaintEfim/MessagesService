package config

import (
	"github.com/joho/godotenv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AuthenticationConfiguration AuthenticationConfiguration `yaml:"AuthenticationConfiguration"`
	EnvironmentVariables        EnvironmentVariables        `yaml:"EnvironmentVariables"`
	Server                      Server                      `yaml:"Server"`
	Redis                       Redis                       `yaml:"Redis"`
	Logs                        Logs                        `yaml:"Logs"`
	Claims                      Claims                      `yaml:"Claims"`
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

type Redis struct {
	Address    string        `yaml:"Address"`
	Password   string        `yaml:"Password"`
	Db         int           `yaml:"Db"`
	Expiration time.Duration `yaml:"Expiration"`
	Timeout    time.Duration `yaml:"Timeout"`
}

type Logs struct {
	Path       string `yaml:"Path"`
	Level      string `yaml:"Level"`
	MaxAge     int    `yaml:"MaxAge"`
	MaxBackups int    `yaml:"MaxBackups"`
}

type Claims struct {
	KeyForId string `yaml:"Key"`
}

func ReadConfig(cfgName, cfgType, cfgPath string) (*Config, error) {
	var cfg Config

	viper.SetConfigName(cfgName)
	viper.SetConfigType(cfgType)
	viper.AddConfigPath(cfgPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg.AuthenticationConfiguration.AccessSecretKey = viper.GetString("ACCESS_SECRET_KEY")

	return &cfg, nil
}
