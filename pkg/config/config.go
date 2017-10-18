package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the configuration file structure
type Config struct {
	Concourse concourse `json:"concourse" mapstructure:"concourse"`
	User      user      `json:"user" mapstructure:"user"`
}

type concourse struct {
	URL     string `json:"url" mapstructure:"url"`
	APIPath string `json:"api_path" mapstructure:"api_path"`

	Teams []team `json:"teams" mapstructure:"teams"`
}

type team struct {
	Name     string `json:"name" mapstructure:"name"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

type user struct {
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

// Load configuration file and override with environment variables
func Load(filename string) *Config {
	viper.SetConfigName(filename)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(fmt.Sprintf("ERROR loading configuration: %s", err))
	}

	// Get env vars as overrides
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	var config Config
	if err = viper.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("ERROR unmarshalling configuration: %s", err))
	}

	return &config
}
