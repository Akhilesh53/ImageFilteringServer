package cloudvisionapi

import (
	"image_filter_server/constants"
	"image_filter_server/src/models"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func AnalyseSafeSearchImage(ctx *gin.Context, imagePath string) (*models.GCPResponse, error) {

	// create a new client
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// open the image file
	f, err := os.Open(imagePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer f.Close()

	// create a new image
	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp, err := client.AnnotateImage(ctx, &visionpb.AnnotateImageRequest{
		Image:    image,
		Features: constants.FeatureList,
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}
	// prepare the response
	gcpSafeSearchResponse := models.DefaultGCPSafeSearch().SetResponse(resp)
	return gcpSafeSearchResponse, nil
}
