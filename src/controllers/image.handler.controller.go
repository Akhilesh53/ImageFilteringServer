package controllers

import (
	apiErr "image_filter_server/pkg/errors"
	"image_filter_server/src/dtos"
	"image_filter_server/src/models"
	"image_filter_server/src/services"
	cloudvisionapi "image_filter_server/src/utils/cloud-vision-api"
	"image_filter_server/src/utils/response"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

	if imageURL == "" {
		response.SendResponse(ctx, nil, apiErr.URLNotPresent.SetUUID(ctx.GetString("uuid")), apiErr.ErrURLNotPresent)
		return
	}

	var docAlreadyPresent bool
	var err error
	var imageFilterApiResponse = dtos.GetDefaultImageFilterAPIResponse()
	var resultString = "image is whitelisted"

	// check if the image url is present in the collection
	if docAlreadyPresent, err = controller.GetImageHandlerService().IsDocPresent(ctx, imageURL); err != nil {
		response.SendResponse(ctx, apiErr.InternalError.SetUUID(ctx.GetString("uuid")), apiErr.InternalError.SetUUID(ctx.GetString("uuid")), errors.WithStack(err))
		return
	}

	if docAlreadyPresent {
		// get the output image url from the collection
		docResp, err := controller.GetImageHandlerService().GetDocResponse(ctx, imageURL)
		if err != nil {
			response.SendResponse(ctx, apiErr.InternalError.SetUUID(ctx.GetString("uuid")), apiErr.InternalError.SetUUID(ctx.GetString("uuid")), errors.WithStack(err))
			return
		}
		imageFilterApiResponse.SetResult(docResp.GetConclusion())
		response.SendResponse(ctx, imageFilterApiResponse, apiErr.RequestProcessSuccess.SetUUID(ctx.GetString("uuid")), nil)
		return
	}

	// hit the google vision api to get the image labels
	googleVisionAPIResp, err := cloudvisionapi.AnalyseSafeSearchImage(ctx, imageURL)

	if err != nil {
		response.SendResponse(ctx, apiErr.InternalError.SetUUID(ctx.GetString("uuid")), apiErr.InternalError.SetUUID(ctx.GetString("uuid")), errors.WithStack(err))
		return
	}

	// check if the image is blacklisted
	CheckSafeSearch(googleVisionAPIResp, &resultString)

	imageFilterApiResponse.SetResult(resultString)

	// store the image response in the collection
	err = controller.GetImageHandlerService().SaveDocResponse(ctx, imageURL, models.NewDefaultFirebaseCollectionResult().SetConclusion(resultString).SetImageURL(imageURL).SetGoogleVisionResult(googleVisionAPIResp))

	if err != nil {
		response.SendResponse(ctx, nil, apiErr.InternalError.SetUUID(ctx.GetString("uuid")), errors.WithStack(err))
		return
	}

	response.SendResponse(ctx, imageFilterApiResponse, apiErr.RequestProcessSuccess.SetUUID(ctx.GetString("uuid")), nil)
}

func CheckSafeSearch(googleVisionAPIResp *models.GCPResponse, resultString *string) {
	*resultString = "image is whitelisted"
	if googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation() != nil {
		// if adult == LIKELY, set string to "image is blacklisted"
		if googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetAdult().String() == "LIKELY" || googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetAdult().String() == "VERY_LIKELY" || googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetAdult().String() == "POSSIBLE" {
			*resultString = "image is blacklisted"
		} else if googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetSpoof().String() == "LIKELY" || googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetSpoof().String() == "VERY_LIKELY" || googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetSpoof().String() == "POSSIBLE" {
			*resultString = "image is blacklisted"
		} else if googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetMedical().String() == "LIKELY" || googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetMedical().String() == "VERY_LIKELY" || googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetMedical().String() == "POSSIBLE" {
			*resultString = "image is blacklisted"
		} else if googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetViolence().String() == "LIKELY" || googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetViolence().String() == "VERY_LIKELY" || googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetViolence().String() == "POSSIBLE" {
			*resultString = "image is blacklisted"
		} else if googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetRacy().String() == "LIKELY" || googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetRacy().String() == "VERY_LIKELY" || googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetRacy().String() == "POSSIBLE" {
			*resultString = "image is blacklisted"
		}
	}
}
