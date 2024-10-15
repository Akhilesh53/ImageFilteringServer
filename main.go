package main

import (
	env "image_filter_server/config"
	"image_filter_server/src/utils/initialisation"

	"github.com/gin-gonic/gin"
)

func main() {

	env.LoadConfig("/Users/b0272559_1/Documents/ImageFilteringServer/")

	// initialise modules
	imageFilterController := initialisation.InitModules()

	// create a new gin router
	router := gin.Default()

	// define routes
	router.POST("/filter", imageFilterController.FilterImage)

}
