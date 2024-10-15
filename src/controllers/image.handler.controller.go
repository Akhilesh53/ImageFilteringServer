package controllers

import (
	"image_filter_server/src/services"
	"sync"

	"github.com/gin-gonic/gin"
)

var imageHandler *ImageHandlerController
var imageHandlerOnce sync.Once

type ImageHandlerController struct {
	imageHandlerService *services.ImageHandlerService
}

// GetImageHandlerController is a function that returns a singleton instance of ImageHandlerController
func GetImageHandlerController() *ImageHandlerController {
	imageHandlerOnce.Do(func() {
		imageHandler = &ImageHandlerController{
			imageHandlerService: services.GetImageHandlerService(),
		}
	})
	return imageHandler
}

func (controller *ImageHandlerController) FilterImage(c *gin.Context) {}
