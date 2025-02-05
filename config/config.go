package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

// App config struct
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   Logger
}

// Server config struct
type ServerConfig struct {
	Port              string
	JwtSecretKey      string
	JwtMaxAge         time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	CtxDefaultTimeout time.Duration
	Debug             bool
}

// Logger config
type Logger struct {
	Level string
}

// Postgresql config
type DatabaseConfig struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	Dbname     string
	DbDriver   string
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigType("yaml")
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
