package constants

import (
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
)

var FeatureList = []*visionpb.Feature{
	// type unspecified
	{
		Type:       visionpb.Feature_TYPE_UNSPECIFIED,
		MaxResults: 1,
	},
	// face detection
	{
		Type:       visionpb.Feature_FACE_DETECTION,
		MaxResults: 1,
	},
	// landmark
	{
		Type:       visionpb.Feature_LANDMARK_DETECTION,
		MaxResults: 1,
	},
	// logo
	{
		Type:       visionpb.Feature_LOGO_DETECTION,
		MaxResults: 1,
	},
	// label
	{
		Type:       visionpb.Feature_LABEL_DETECTION,
		MaxResults: 1,
	},
	// text
	{
		Type:       visionpb.Feature_TEXT_DETECTION,
		MaxResults: 1,
	},
	//doc text
	{
		Type:       visionpb.Feature_DOCUMENT_TEXT_DETECTION,
		MaxResults: 1,
	},
	// safe search
	{
		Type:       visionpb.Feature_SAFE_SEARCH_DETECTION,
		MaxResults: 1,
	},
	// image properties
	{
		Type:       visionpb.Feature_IMAGE_PROPERTIES,
		MaxResults: 1,
	},
	// crop hints
	{
		Type:       visionpb.Feature_CROP_HINTS,
		MaxResults: 1,
	},
	// web detection
	{
		Type:       visionpb.Feature_WEB_DETECTION,
		MaxResults: 1,
	},
	// product search
	{
		Type:       visionpb.Feature_PRODUCT_SEARCH,
		MaxResults: 1,
	},
	// object localization
	{
		Type:       visionpb.Feature_OBJECT_LOCALIZATION,
		MaxResults: 1,
	},
}
