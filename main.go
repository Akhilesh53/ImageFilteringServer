package main

import (
	"image_filter_server/pkg/logging"
	"image_filter_server/src/utils/initialisation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	// create a new gin router
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	router := gin.New()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = true

	// initialise modules
	imageFilterController := initialisation.InitModules()
	
	// define routes
	router.POST("/filter", imageFilterController.FilterImage)

	// run the server by printing
	server := &http.Server{
		Addr:    ":" + viper.GetString("PORT"),
		Handler: router,
	}

	logging.Info(&gin.Context{}, "Server is running on port "+viper.GetString("PORT"))

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
