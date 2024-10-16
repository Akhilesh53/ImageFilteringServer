package cloudvisionapi

import (
	"image_filter_server/src/models"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func AnalyseSafeSearchImage(ctx *gin.Context, imagePath string) (*models.GCPSafeSearch, error) {

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

	// detect safe search properties
	props, err := client.DetectSafeSearch(ctx, image, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// prepare the response
	gcpSafeSearchResponse := models.DefaultGCPSafeSearch().SetAdult(props.GetAdult().String()).SetMedical(props.GetMedical().String()).SetSpoof(props.GetSpoof().String()).SetViolence(props.GetViolence().String()).SetRacy(props.GetRacy().String())

	return gcpSafeSearchResponse, nil
}
