package controllers

import (
	"fmt"
	apiErr "image_filter_server/pkg/errors"
	"image_filter_server/src/dtos"
	"image_filter_server/src/models"
	"image_filter_server/src/services"
	cloudvisionapi "image_filter_server/src/utils/cloud-vision-api"
	"image_filter_server/src/utils/response"
	"strings"
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
		fmt.Println("image url already present in the collection")
		imageFilterApiResponse.SetResult(docResp.GetConclusion())
		response.SendResponse(ctx, imageFilterApiResponse, apiErr.RequestProcessSuccess.SetUUID(ctx.GetString("uuid")), nil)
		return
	}

	// hit the google vision api to get the image labels
	googleVisionAPIResp, err := cloudvisionapi.AnnotateImage(ctx, imageURL)

	if err != nil {
		response.SendResponse(ctx, apiErr.InternalError.SetUUID(ctx.GetString("uuid")), apiErr.InternalError.SetUUID(ctx.GetString("uuid")), errors.WithStack(err))
		return
	}

	if blockedWordsPresent := controller.GetImageHandlerService().IsBlockedWordsPresent(ctx, collectWordsFromResponse(googleVisionAPIResp)); blockedWordsPresent {
		resultString = "image is blacklisted"
	}
	imageFilterApiResponse.SetResult(resultString)

	// store the image response in the collection
	err = controller.GetImageHandlerService().SaveDocResponse(ctx, imageURL, models.NewDefaultFirebaseCollectionResult().SetConclusion(resultString).SetImageURL(imageURL).SetGoogleVisionResult(googleVisionAPIResp))

	if err != nil {
		response.SendResponse(ctx, nil, apiErr.InternalError.SetUUID(ctx.GetString("uuid")), errors.WithStack(err))
		return
	}

	response.SendResponse(ctx, imageFilterApiResponse, apiErr.RequestProcessSuccess.SetUUID(ctx.GetString("uuid")), nil)
}

func collectWordsFromResponse(googleVisionAPIResp *models.GCPResponse) []string {

	var wg sync.WaitGroup

	var (
		wordsChan = make(chan []string, 9)
		words     = []string{}
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		getFaceAnnotationsWords(googleVisionAPIResp, wordsChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		landmarkAnnotationsWords(googleVisionAPIResp, wordsChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logoAnnotationsWords(googleVisionAPIResp, wordsChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		labelAnnotationsWords(googleVisionAPIResp, wordsChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		localizedObjectAnnotationsWords(googleVisionAPIResp, wordsChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		textAnnotationsWords(googleVisionAPIResp, wordsChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fullTextAnnotationsWords(googleVisionAPIResp, wordsChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		safeSearchAnnotationsWords(googleVisionAPIResp, wordsChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		productSearchAnnotationsWords(googleVisionAPIResp, wordsChan)
	}()

	wg.Wait()
	close(wordsChan)

	for wordArr := range wordsChan {
		if len(wordArr) > 0 {
			words = append(words, wordArr...)
		}
	}
	return words
}

func getFaceAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	// face annotations
	if googleVisionAPIResp.GetResponse().GetFaceAnnotations() != nil {
		for _, faceAnnotation := range googleVisionAPIResp.GetResponse().GetFaceAnnotations() {
			if faceAnnotation.GetAngerLikelihood().String() == "LIKELY" ||
				faceAnnotation.GetAngerLikelihood().String() == "VERY_LIKELY" ||
				faceAnnotation.GetAngerLikelihood().String() == "POSSIBLE" {
				words = append(words, "anger")
			}
		}
	}
	wordsChan <- words
}

// land mark annotations
func landmarkAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	// landmark annotations
	if googleVisionAPIResp.GetResponse().GetLandmarkAnnotations() != nil {
		for _, landmarkAnnotation := range googleVisionAPIResp.GetResponse().GetLandmarkAnnotations() {
			words = append(words, strings.ToLower(landmarkAnnotation.GetDescription()))

			// todo: check whhether to consider name or values in properties
			// properties
			if landmarkAnnotation.GetProperties() != nil {
				for _, property := range landmarkAnnotation.GetProperties() {
					words = append(words, strings.ToLower(property.GetName()))
				}
			}
		}
	}
	wordsChan <- words
}

// logo annotations
func logoAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	// logo annotations
	if googleVisionAPIResp.GetResponse().GetLogoAnnotations() != nil {
		for _, logoAnnotation := range googleVisionAPIResp.GetResponse().GetLogoAnnotations() {
			words = append(words, strings.ToLower(logoAnnotation.GetDescription()))

			// properties
			if logoAnnotation.GetProperties() != nil {
				for _, property := range logoAnnotation.GetProperties() {
					words = append(words, strings.ToLower(property.GetName()))
				}
			}
		}
	}
	wordsChan <- words
}

// label
func labelAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	// label annotations
	if googleVisionAPIResp.GetResponse().GetLabelAnnotations() != nil {
		for _, labelAnnotation := range googleVisionAPIResp.GetResponse().GetLabelAnnotations() {
			words = append(words, strings.ToLower(labelAnnotation.GetDescription()))

			// properties
			if labelAnnotation.GetProperties() != nil {
				for _, property := range labelAnnotation.GetProperties() {
					words = append(words, strings.ToLower(property.GetName()))
				}
			}
		}
	}
	wordsChan <- words
}

// localised object annotations
func localizedObjectAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	// localized object annotations
	if googleVisionAPIResp.GetResponse().GetLocalizedObjectAnnotations() != nil {
		for _, localizedObjectAnnotation := range googleVisionAPIResp.GetResponse().GetLocalizedObjectAnnotations() {
			words = append(words, strings.ToLower(localizedObjectAnnotation.GetName()))
		}
	}
	wordsChan <- words
}

// text annotations
func textAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	// text annotations
	if googleVisionAPIResp.GetResponse().GetTextAnnotations() != nil {
		for _, textAnnotation := range googleVisionAPIResp.GetResponse().GetTextAnnotations() {
			words = append(words, strings.ToLower(textAnnotation.GetDescription()))

			// properties
			if textAnnotation.GetProperties() != nil {
				for _, property := range textAnnotation.GetProperties() {
					words = append(words, strings.ToLower(property.GetName()))
				}
			}
		}
	}
	wordsChan <- words
}

// full text annotations
func fullTextAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	// full text annotations
	if googleVisionAPIResp.GetResponse().GetFullTextAnnotation() != nil {
		words = append(words, strings.ToLower(googleVisionAPIResp.GetResponse().GetFullTextAnnotation().GetText()))
	}
	wordsChan <- words
}

// safe search annotations
func safeSearchAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	// safe search annotations
	if googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation() != nil {
		if googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetAdult().String() == "LIKELY" ||
			googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetAdult().String() == "VERY_LIKELY" ||
			googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetAdult().String() == "POSSIBLE" {
			words = append(words, "adult")
		}

		if googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetSpoof().String() == "LIKELY" ||
			googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetSpoof().String() == "VERY_LIKELY" ||
			googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetSpoof().String() == "POSSIBLE" {
			words = append(words, "spoof")
		}

		if googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetMedical().String() == "LIKELY" ||
			googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetMedical().String() == "VERY_LIKELY" ||
			googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetMedical().String() == "POSSIBLE" {
			words = append(words, "medical")
		}

		if googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetViolence().String() == "LIKELY" ||
			googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetViolence().String() == "VERY_LIKELY" ||
			googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetViolence().String() == "POSSIBLE" {
			words = append(words, "violence")
		}

		if googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetRacy().String() == "LIKELY" ||
			googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetRacy().String() == "VERY_LIKELY" ||
			googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation().GetRacy().String() == "POSSIBLE" {
			words = append(words, "racy")
		}
	}
	wordsChan <- words
}

// product search annotations
func productSearchAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	// product search annotations
	if googleVisionAPIResp.GetResponse().GetProductSearchResults() != nil {
		// product read results: grouped
		if googleVisionAPIResp.GetResponse().GetProductSearchResults().GetProductGroupedResults() != nil {
			for _, productGroupedResult := range googleVisionAPIResp.GetResponse().GetProductSearchResults().GetProductGroupedResults() {
				for _, productResult := range productGroupedResult.GetResults() {
					words = append(words, strings.ToLower(productResult.GetProduct().GetDisplayName()))
				}
			}
		}

		// product read results: object
		if googleVisionAPIResp.GetResponse().GetProductSearchResults().GetResults() != nil {
			for _, productResult := range googleVisionAPIResp.GetResponse().GetProductSearchResults().GetResults() {
				words = append(words, strings.ToLower(productResult.GetProduct().GetDisplayName()))
			}
		}
	}
	wordsChan <- words
}
