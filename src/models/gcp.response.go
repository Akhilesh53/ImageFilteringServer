package models

import (
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
)

type GCPResponse struct {
	Response *visionpb.AnnotateImageResponse `json:"responses"`
}

// func to create new instance of GCPSafeSearch
func NewGCPSafeSearch(response *visionpb.AnnotateImageResponse) *GCPResponse {
	return &GCPResponse{
		Response: response,
	}
}

// func to get default instance of GCPSafeSearch
func DefaultGCPSafeSearch() *GCPResponse {
	return &GCPResponse{
	}
}

// func to get response
func (g *GCPResponse) GetResponse() *visionpb.AnnotateImageResponse {
	return g.Response
}

// func to set response
func (g *GCPResponse) SetResponse(response *visionpb.AnnotateImageResponse) *GCPResponse {
	g.Response = response
	return g
}

