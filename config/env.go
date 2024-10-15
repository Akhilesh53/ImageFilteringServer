package env

import "github.com/spf13/viper"

// Load configuration from environment variables
func LoadConfig(configPath string) {

	// Set the configuration path
	viper.AddConfigPath(configPath)

	// Set the configuration file name
	viper.SetConfigName("imagefilter.env")

	// set the configuration type
	viper.SetConfigType("env")

	// Read the configuration file
	viper.AutomaticEnv()
	if err:= viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

