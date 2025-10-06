package config

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	SubServiceConfig struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"sub-service-config"`

	DatabaseConfig struct {
		Url string `mapstructure:"url"`
	} `mapstructure:"sub-database-config"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("read config file: %w", err)
	}

	expanded := os.ExpandEnv(string(data))

	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(strings.NewReader(expanded)); err != nil {
		return config, fmt.Errorf("viper read config: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("unmarshal config: %w", err)
	}

	return config, nil
}
