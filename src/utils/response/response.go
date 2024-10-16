package response

import (
	apiErr "image_filter_server/pkg/errors"
	"image_filter_server/pkg/logging"
	"image_filter_server/src/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SendResponse(ctx *gin.Context, payload interface{}, Err *apiErr.Error, err error) {
	if payload == nil {
		payload = Err
		logging.Error(ctx, Err.ErrorMessage, zap.Error(err))
	}
	if err == nil {
		logging.Info(ctx, Err.ErrorMessage)
	}
	logging.Info(ctx, "Response", zap.Any("data", payload))
	ctx.Writer.Write(utils.InterfaceToBytes(ctx, payload))
	ctx.AbortWithStatus(Err.StatusCode)
}
