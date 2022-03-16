package util

import (
	"os"

	"github.com/spf13/viper"
)

//Config stores all configuration fields
type Config struct {
	DBSource string `mapstructure:"DB_SOURCE"`
}

//Loads configuration from the file at the provided path
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.SetEnvPrefix("DB_SOURCE")
	//overwrite config with environment variables if they exist
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		if s := os.Getenv("DB_SOURCE"); s != "" {
			config.DBSource = s
			err = nil
		}
		return
	}
	err = viper.Unmarshal(&config)
	return
}
