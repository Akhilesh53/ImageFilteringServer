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
	}
	if Err != nil {
		logging.Info(ctx, Err.ErrorMessage)
	}
	if err != nil {
		logging.Error(ctx, "error : ", zap.Error(err))
	}

	logging.Info(ctx, "Response", zap.Any("data", payload))
	ctx.Writer.Write(utils.InterfaceToBytes(ctx, payload))
	ctx.AbortWithStatus(Err.StatusCode)
}
