package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	DBDriver              string        `mapstructure:"DB_DRIVER"`
	DBSource              string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress     string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress     string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey     string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration   time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	ResfreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Don't fail if the config file is missing — for CI environments
	if err = viper.ReadInConfig(); err != nil {
		fmt.Println("Config file not found, using only environment variables")
	}

	err = viper.Unmarshal(&config)
	return
}
