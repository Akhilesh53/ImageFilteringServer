package services

import (
	"image_filter_server/src/daos"
	"image_filter_server/src/models"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

// get image handler dao
func (service *ImageHandlerService) GetImageHandlerDao() *daos.ImageHandlerDao {
	return service.imageHandlerDao
}

// servuce func to check image url is already present or not
func (service *ImageHandlerService) IsImageURLPresent(ctx *gin.Context, imageURL string) (bool, error) {
	return service.GetImageHandlerDao().IsImageURLPresent(ctx, imageURL)
}

// get image response from colllection
func (service *ImageHandlerService) GetImageUrlResponse(ctx *gin.Context, imageURL string) (*models.FirebaseCollectionResult, error) {
	resp, err := service.GetImageHandlerDao().GetImageUrlResponse(ctx, imageURL)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	// todo:
	// prepare the response struct and share it

	return resp, nil
}
