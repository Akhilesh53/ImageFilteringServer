package main

import (
	"fmt"
	env "image_filter_server/config"

	"github.com/spf13/viper"
)

func main() {

	env.LoadConfig("/Users/b0272559_1/Documents/ImageFilteringServer/")
	// print all env variables using viper

	fmt.Println(viper.GetString("LOG_LEVEL"))
}
