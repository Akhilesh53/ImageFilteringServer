package dtos

// request
type ImageFilterAPIRequest struct {
	ImageUrl string `json:"image_url"`
}

// func to get the image filter api request
func GetImageFilterAPIRequest(imageUrl string) *ImageFilterAPIRequest {
	return &ImageFilterAPIRequest{
		ImageUrl: imageUrl,
	}
}

// get default image filter api request
func GetDefaultImageFilterAPIRequest() *ImageFilterAPIRequest {
	return &ImageFilterAPIRequest{}
}

// set image url
func (request *ImageFilterAPIRequest) SetImageUrl(imageUrl string) *ImageFilterAPIRequest {
	request.ImageUrl = imageUrl
	return request
}

// get image url
func (request *ImageFilterAPIRequest) GetImageUrl() string {
	return request.ImageUrl
}
