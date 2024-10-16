package utils

import (
	"encoding/json"
	"image_filter_server/pkg/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InterfaceToBytes(ctx *gin.Context, v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		logging.Error(ctx, "Error converting interface to bytes", zap.Any("error", err))
	}
	return b
}
