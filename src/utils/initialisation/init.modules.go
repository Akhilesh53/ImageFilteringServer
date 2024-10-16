package initialisation

import (
	env "image_filter_server/config"
	"image_filter_server/pkg/logging"
	"image_filter_server/src/controllers"
)

func InitModules(envPath string) *controllers.ImageHandlerController {
	env.LoadConfig(envPath)
	logging.InitialiseLogger()
	imageFilterController := controllers.GetImageHandlerController()
	return imageFilterController
}
