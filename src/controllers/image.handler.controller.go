package controllers

import (
	apiErr "image_filter_server/pkg/errors"
	"image_filter_server/src/services"
	cloudvisionapi "image_filter_server/src/utils/cloud-vision-api"
	"image_filter_server/src/utils/response"
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

// get image service
func (controller *ImageHandlerController) GetImageHandlerService() *services.ImageHandlerService {
	return controller.imageHandlerService
}

func (controller *ImageHandlerController) FilterImage(ctx *gin.Context) {
	// get the image url from the request query
	imageURL := ctx.Query("url")

	var urlAlreadyPresent bool
	var err error
	// check if the image url is present in the collection
	if urlAlreadyPresent, err = controller.GetImageHandlerService().IsImageURLPresent(ctx, imageURL); err != nil {
		response.SendResponse(ctx, apiErr.InternalError.SetUUID(ctx.GetString("uuid")), apiErr.InternalError.SetUUID(ctx.GetString("uuid")), err)
		return
	}

	if urlAlreadyPresent {
		// get the output image url from the collection
		urlResponse, err := controller.GetImageHandlerService().GetImageUrlResponse(ctx, imageURL)
		if err != nil {
			response.SendResponse(ctx, apiErr.InternalError.SetUUID(ctx.GetString("uuid")), apiErr.InternalError.SetUUID(ctx.GetString("uuid")), err)
			return
		}
		response.SendResponse(ctx, urlResponse, apiErr.InternalError.SetUUID(ctx.GetString("uuid")), nil)
		return
	}

	// hit the google vision api to get the image labels
	safeSearchResp, err := cloudvisionapi.AnalyseSafeSearchImage(ctx, imageURL)

	if err != nil {
		response.SendResponse(ctx, apiErr.InternalError.SetUUID(ctx.GetString("uuid")), apiErr.InternalError.SetUUID(ctx.GetString("uuid")), err)
		return
	}

	response.SendResponse(ctx, safeSearchResp, nil, nil)

	// if api doesnot extract anything , means len(result) == 0, image is safe
	// else read triggering words from json :

	/*
		if result["adult"] == "LIKELY" || result["violence"] == "LIKELY" || result["racy"] == "LIKELY" ||
			result["adult"] == "VERY_LIKELY" || result["violence"] == "VERY_LIKELY" || result["racy"] == "VERY_LIKELY" {
			conclusion = "blocked"
		}
	*/

	// store the result in file store
}
