package initialisation

import (
	"image_filter_server/pkg/logging"
	"image_filter_server/src/controllers"
)

func InitModules() *controllers.ImageHandlerController {
	logging.InitialiseLogger()
	imageFilterController := controllers.GetImageHandlerController()
	return imageFilterController
}
