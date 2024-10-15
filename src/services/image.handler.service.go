package services

import (
	"image_filter_server/src/daos"
	"sync"
)

var imageHandlerService *ImageHandlerService
var imageHandlerServiceOnce sync.Once

type ImageHandlerService struct {
	imageHandlerDao *daos.ImageHandlerDao
}

// GetImageHandlerService is a function that returns a singleton instance of ImageHandlerService
func GetImageHandlerService() *ImageHandlerService {
	imageHandlerServiceOnce.Do(func() {
		imageHandlerService = &ImageHandlerService{
			imageHandlerDao: daos.GetImageHandlerDao(),
		}
	})
	return imageHandlerService
}
