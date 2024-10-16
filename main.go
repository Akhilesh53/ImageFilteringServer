package main

import (
	"fmt"
	"image_filter_server/src/utils/initialisation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// initialise modules
	imageFilterController := initialisation.InitModules("/Users/akhileshmahajan/Documents/ImageFilteringServer")

	// create a new gin router
	router := gin.Default()

	// define routes
	router.POST("/filter", imageFilterController.FilterImage)

	// run the server by printing
	server := &http.Server{
		Addr:    ":" + viper.GetString("PORT"),
		Handler: router,
	}

	fmt.Println("Starting server on port: " + viper.GetString("PORT"))

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
