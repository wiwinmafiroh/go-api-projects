package helpers

import (
	"06-go-api-with-unittest/pkg/errs"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetParamId(ctx *gin.Context, key string) (uint, errs.ErrorMessage) {
	value := ctx.Param(key)

	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, errs.NewBadRequestError("Parameter id must be a valid integer")
	}

	if id < 0 {
		return 0, errs.NewBadRequestError("Parameter id must be a non-negative integer")
	}

	return uint(id), nil
}
