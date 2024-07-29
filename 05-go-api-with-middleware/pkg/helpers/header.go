package helpers

import (
	"05-go-api-with-middleware/pkg/errs"

	"github.com/gin-gonic/gin"
)

func CheckContentType(ctx *gin.Context) errs.ErrorMessage {
	if ctx.GetHeader("Content-Type") != "application/json" {
		return errs.NewUnsupportedMediaTypeError("Unsupported content-type. Only 'application/json' is supported.")
	}

	return nil
}
