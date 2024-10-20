package cloudvisionapi

import (
	"context"
	"fmt"
	"image_filter_server/constants"
	"image_filter_server/src/models"
	"net/http"
	"os"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Create a new Vision API client
func createVisionClient(ctx context.Context) (*vision.ImageAnnotatorClient, error) {
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client, nil
}

// Load image from either local path or URL
func loadImage(imagePath string) (*visionpb.Image, error) {
	if isURL(imagePath) {
		// Load image from URL
		return loadImageFromURL(imagePath)
	}
	// Load image from local file system
	return loadImageFromFile(imagePath)
}

// Helper function to check if the imagePath is a URL
func isURL(imagePath string) bool {
	return strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://")
}

// Load image from a URL
func loadImageFromURL(imageURL string) (*visionpb.Image, error) {
	fmt.Println("Loading image from URL:", imageURL)
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("failed to fetch image from URL: %s", resp.Status)
	}

	imageData, err := vision.NewImageFromReader(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return imageData, nil
}

// Load image from local file system
func loadImageFromFile(filePath string) (*visionpb.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer f.Close()

	imageData, err := vision.NewImageFromReader(f)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return imageData, nil
}

// Annotate the image using Vision API
func annotateImage(client *vision.ImageAnnotatorClient, ctx context.Context, image *visionpb.Image) (*models.GCPResponse, error) {
	resp, err := client.AnnotateImage(ctx, &visionpb.AnnotateImageRequest{
		Image:    image,
		Features: constants.FeatureList, // Ensure this contains the required features for analysis
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}

	// Prepare the response
	gcpSafeSearchResponse := models.DefaultGCPSafeSearch().SetResponse(resp)
	return gcpSafeSearchResponse, nil
}

// Main function to annotate an image from either a local path or a URL
func AnnotateImage(ctx *gin.Context, imagePath string) (*models.GCPResponse, error) {
	// Create the Vision client
	client, err := createVisionClient(ctx)
	if err != nil {
		return nil, err
	}

	// Load the image (either from local path or URL)
	image, err := loadImage(imagePath)
	if err != nil {
		return nil, err
	}

	// Annotate the image
	return annotateImage(client, ctx, image)
}
