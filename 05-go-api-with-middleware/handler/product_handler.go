package handler

import (
	"05-go-api-with-middleware/dto"
	"05-go-api-with-middleware/entity"
	"05-go-api-with-middleware/pkg/errs"
	"05-go-api-with-middleware/pkg/helpers"
	"05-go-api-with-middleware/service"

	"github.com/gin-gonic/gin"
)

type productHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) productHandler {
	return productHandler{
		productService: productService,
	}
}

func (p *productHandler) CreateProduct(ctx *gin.Context) {
	err := helpers.CheckContentType(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	var productRequest dto.ProductRequest

	if err := ctx.ShouldBindJSON(&productRequest); err != nil {
		errBindJSON := errs.NewUnprocessableEntityError("Invalid request body")

		ctx.AbortWithStatusJSON(errBindJSON.StatusCode(), errBindJSON)
		return
	}

	userID := ctx.MustGet("userData").(entity.User).Id

	result, err := p.productService.CreateProduct(userID, productRequest)
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	ctx.JSON(result.StatusCode, result)
}

func (p *productHandler) GetProducts(ctx *gin.Context) {
	userData := ctx.MustGet("userData").(entity.User)

	result, err := p.productService.GetProducts(userData.Id, string(userData.Role))
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	ctx.JSON(result.StatusCode, result)
}

func (p *productHandler) GetProductById(ctx *gin.Context) {
	productId, err := helpers.GetParamId(ctx, "productId")
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	result, err := p.productService.GetProductById(productId)
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	ctx.JSON(result.StatusCode, result)
}

func (p *productHandler) UpdateProductById(ctx *gin.Context) {
	productId, err := helpers.GetParamId(ctx, "productId")
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	err = helpers.CheckContentType(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	var productRequest dto.ProductRequest

	if err := ctx.ShouldBindJSON(&productRequest); err != nil {
		errBindJSON := errs.NewUnprocessableEntityError("Invalid request body")

		ctx.AbortWithStatusJSON(errBindJSON.StatusCode(), errBindJSON)
		return
	}

	result, err := p.productService.UpdateProductById(productId, productRequest)
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	ctx.JSON(result.StatusCode, result)
}

func (p *productHandler) DeleteProductById(ctx *gin.Context) {
	productId, err := helpers.GetParamId(ctx, "productId")
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	result, err := p.productService.DeleteProductById(productId)
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	ctx.JSON(result.StatusCode, result)
}
