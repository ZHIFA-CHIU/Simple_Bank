package utils

import "github.com/spf13/viper"

// Store configuration settings for the application.
type Config struct {
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBSource     string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig reads configuration from a file at the specified path and returns a Config struct.
func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv() // read environment variables that match the config keys
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}