package dtos

type ImageFilterAPIResponse struct {
	Result string `json:"result"`
}

// func to get the image filter api response
func GetImageFilterAPIResponse(result string) *ImageFilterAPIResponse {
	return &ImageFilterAPIResponse{
		Result: result,
	}
}

// get default image filter api response
func GetDefaultImageFilterAPIResponse() *ImageFilterAPIResponse {
	return &ImageFilterAPIResponse{}
}

// set result
func (response *ImageFilterAPIResponse) SetResult(result string) *ImageFilterAPIResponse {
	response.Result = result
	return response
}

// get result
func (response *ImageFilterAPIResponse) GetResult() string {
	return response.Result
}
