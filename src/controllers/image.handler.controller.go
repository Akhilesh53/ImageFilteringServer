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

// get image service
func (controller *ImageHandlerController) GetImageHandlerService() *services.ImageHandlerService {
	return controller.imageHandlerService
}

func (controller *ImageHandlerController) FilterImage(c *gin.Context) {
	// get the image url from the request query
	imageURL := c.Query("url")

	var urlAlreadyPresent bool
	var err error
	// check if the image url is present in the collection
	if urlAlreadyPresent, err = controller.GetImageHandlerService().CheckImageURLPresent(imageURL); err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	if urlAlreadyPresent {
		// get the output image url from the collection
		// return the output image url
	}

	// hit the google vision api to get the image labels
	/*
		result := map[string]interface{}{
			"adult":     resp.Adult.String(),
			"violence":  resp.Violence.String(),
			"medical":   resp.Medical.String(),
			"racy":      resp.Racy.String(),
			"spoof":     resp.Spoof.String(),
		}
	*/

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
