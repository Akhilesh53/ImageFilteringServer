package controllers

import (
	apiErr "image_filter_server/pkg/errors"
	"image_filter_server/pkg/logging"
	"image_filter_server/src/dtos"
	"image_filter_server/src/models"
	"image_filter_server/src/services"
	cloudvisionapi "image_filter_server/src/utils/cloud-vision-api"
	"image_filter_server/src/utils/response"
	"strings"
	"sync"

	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"go.uber.org/zap"

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

func (controller *ImageHandlerController) VerifyImage(ctx *gin.Context) {

	// bind the json request to the struct
	var request dtos.ImageFilterAPIRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.SendResponse(ctx, nil, apiErr.InvalidRequest.SetUUID(ctx.GetString("uuid")), errors.WithStack(err))
		return
	}

	// get the image url from the request query
	imageURL := request.GetImageUrl()
	logging.Info(ctx, "image url received : ", zap.String("image url : ", imageURL))

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
		logging.Debug(ctx, "doc response already present ", zap.String(" url : ", docResp.GetImageURL()))
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

	blockedWords := collectWordsFromResponse(googleVisionAPIResp)
	if blockedWordsPresent := controller.GetImageHandlerService().IsBlockedWordsPresent(ctx, blockedWords); blockedWordsPresent {
		logging.Debug(ctx, " blocked words : ", zap.Strings(" blocked words : ", blockedWords))
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

// Main function to collect words from Google Vision API response
func collectWordsFromResponse(googleVisionAPIResp *models.GCPResponse) []string {
	var wg sync.WaitGroup
	wordsChan := make(chan []string, 9)
	words := []string{}

	// Collect words in parallel from different annotations
	annotationFuncs := []func(*models.GCPResponse, chan []string){
		getFaceAnnotationsWords,
		landmarkAnnotationsWords,
		logoAnnotationsWords,
		labelAnnotationsWords,
		localizedObjectAnnotationsWords,
		textAnnotationsWords,
		fullTextAnnotationsWords,
		safeSearchAnnotationsWords,
		productSearchAnnotationsWords,
	}

	for _, annotationFunc := range annotationFuncs {
		wg.Add(1)
		go func(annotateFunc func(*models.GCPResponse, chan []string)) {
			defer wg.Done()
			annotateFunc(googleVisionAPIResp, wordsChan)
		}(annotationFunc)
	}

	// Wait for all routines to complete
	wg.Wait()
	close(wordsChan)

	// Collect all words from channels
	for wordArr := range wordsChan {
		if len(wordArr) > 0 {
			words = append(words, wordArr...)
		}
	}

	return words
}

// Face Annotations
func getFaceAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	if annotations := googleVisionAPIResp.GetResponse().GetFaceAnnotations(); annotations != nil {
		for _, faceAnnotation := range annotations {
			if isAngerLikely(faceAnnotation.GetAngerLikelihood().String()) {
				words = append(words, "anger")
			}
		}
	}
	wordsChan <- words
}

func isAngerLikely(likelihood string) bool {
	return likelihood == "LIKELY" || likelihood == "VERY_LIKELY" || likelihood == "POSSIBLE"
}

// Landmark Annotations
func landmarkAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	if annotations := googleVisionAPIResp.GetResponse().GetLandmarkAnnotations(); annotations != nil {
		for _, landmark := range annotations {
			words = append(words, strings.ToLower(landmark.GetDescription()))
			words = append(words, collectProperties(landmark.GetProperties())...)
		}
	}
	wordsChan <- words
}

// Logo Annotations
func logoAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	if annotations := googleVisionAPIResp.GetResponse().GetLogoAnnotations(); annotations != nil {
		for _, logo := range annotations {
			words = append(words, strings.ToLower(logo.GetDescription()))
			words = append(words, collectProperties(logo.GetProperties())...)
		}
	}
	wordsChan <- words
}

// Label Annotations
func labelAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	if annotations := googleVisionAPIResp.GetResponse().GetLabelAnnotations(); annotations != nil {
		for _, label := range annotations {
			words = append(words, strings.ToLower(label.GetDescription()))
			words = append(words, collectProperties(label.GetProperties())...)
		}
	}
	wordsChan <- words
}

// Localized Object Annotations
func localizedObjectAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	if annotations := googleVisionAPIResp.GetResponse().GetLocalizedObjectAnnotations(); annotations != nil {
		for _, obj := range annotations {
			words = append(words, strings.ToLower(obj.GetName()))
		}
	}
	wordsChan <- words
}

// Text Annotations
func textAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	if annotations := googleVisionAPIResp.GetResponse().GetTextAnnotations(); annotations != nil {
		for _, text := range annotations {
			words = append(words, strings.ToLower(text.GetDescription()))
			words = append(words, collectProperties(text.GetProperties())...)
		}
	}
	wordsChan <- words
}

// Full Text Annotations
func fullTextAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	if fullText := googleVisionAPIResp.GetResponse().GetFullTextAnnotation(); fullText != nil {
		words = append(words, strings.ToLower(fullText.GetText()))
	}
	wordsChan <- words
}

// Safe Search Annotations
func safeSearchAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	if safeSearch := googleVisionAPIResp.GetResponse().GetSafeSearchAnnotation(); safeSearch != nil {
		words = append(words, safeSearchChecks(safeSearch)...)
	}
	wordsChan <- words
}

func safeSearchChecks(safeSearch *visionpb.SafeSearchAnnotation) []string {
	var words []string
	if isLikely(safeSearch.GetAdult()) {
		words = append(words, "adult")
	}
	if isLikely(safeSearch.GetSpoof()) {
		words = append(words, "spoof")
	}
	if isLikely(safeSearch.GetMedical()) {
		words = append(words, "medical")
	}
	if isLikely(safeSearch.GetViolence()) {
		words = append(words, "violence")
	}
	if isLikely(safeSearch.GetRacy()) {
		words = append(words, "racy")
	}
	return words
}

func isLikely(likelihood visionpb.Likelihood) bool {
	return likelihood.String() == "LIKELY" || likelihood.String() == "VERY_LIKELY" || likelihood.String() == "POSSIBLE"
}

// Product Search Annotations
func productSearchAnnotationsWords(googleVisionAPIResp *models.GCPResponse, wordsChan chan []string) {
	var words []string
	if productSearch := googleVisionAPIResp.GetResponse().GetProductSearchResults(); productSearch != nil {
		words = append(words, collectProductSearchWords(productSearch)...)
	}
	wordsChan <- words
}

func collectProductSearchWords(productSearch *visionpb.ProductSearchResults) []string {
	var words []string
	if groupedResults := productSearch.GetProductGroupedResults(); groupedResults != nil {
		for _, group := range groupedResults {
			for _, result := range group.GetResults() {
				words = append(words, strings.ToLower(result.GetProduct().GetDisplayName()))
			}
		}
	}
	if results := productSearch.GetResults(); results != nil {
		for _, result := range results {
			words = append(words, strings.ToLower(result.GetProduct().GetDisplayName()))
		}
	}
	return words
}

// Helper function to collect properties from annotation
func collectProperties(properties []*visionpb.Property) []string {
	var words []string
	for _, property := range properties {
		words = append(words, strings.ToLower(property.GetName()))
	}
	return words
}
