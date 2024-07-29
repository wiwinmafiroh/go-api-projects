package helpers

import (
	"06-go-api-with-unittest/pkg/errs"

	"github.com/gin-gonic/gin"
)

func CheckContentType(ctx *gin.Context) errs.ErrorMessage {
	if ctx.GetHeader("Content-Type") != "application/json" {
		return errs.NewUnsupportedMediaTypeError("Unsupported content-type. Only 'application/json' is supported.")
	}

	return nil
}
