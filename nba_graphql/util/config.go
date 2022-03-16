package util

import (
	"fmt"

	"github.com/sirupsen/logrus"
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
	logrus.Warn("read")
	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			v := viper.Get("DB_SOURCE")
			if v == nil {
				err = fmt.Errorf("config file not found and db source not set")
				return
			}
		} else {
			// Config file was found but another error was produced
			return
		}
	}
	err = viper.Unmarshal(&config)
	return
}
