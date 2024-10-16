package utils

import (
	"encoding/json"
	"image_filter_server/pkg/logging"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var reg = regexp.MustCompile("[^a-zA-Z0-9]+")

func InterfaceToBytes(ctx *gin.Context, v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		logging.Error(ctx, "Error converting interface to bytes", zap.Any("error", err))
	}
	return b
}

func ModifyURL(url string) string {
	// all unnecessary characters are removed from the url replaced by _
	url = reg.ReplaceAllString(url, "_")
	return url
}
