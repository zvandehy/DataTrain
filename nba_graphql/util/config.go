package util

import "github.com/spf13/viper"

//Config stores all configuration fields
type Config struct {
	DBSource string `mapstructure:"DB_SOURCE"`
}

//Loads configuration from the file at the provided path
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	//overwrite config with environment variables if they exist
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
