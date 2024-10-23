package initialisation

import (
	env "image_filter_server/config"
	apiErr "image_filter_server/pkg/errors"
	"image_filter_server/pkg/logging"
	"image_filter_server/src/controllers"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitModules() *controllers.ImageHandlerController {

	// get current working directory using library func
	envPath, err := os.Getwd()
	if err != nil {
		logging.Error(&gin.Context{}, apiErr.ErrEnvFilePathNotFound.Error(), zap.Error(err))
		panic(err)
	}
	// load
	env.LoadConfig(envPath)
	logging.InitialiseLogger()
	imageFilterController := controllers.GetImageHandlerController()
	return imageFilterController
}
