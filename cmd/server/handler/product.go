package handler

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/matias-ziliotto/HackthonGo/internal/product"
	"github.com/matias-ziliotto/HackthonGo/pkg/web"
)

type ProductHandler struct {
	productService product.ProductService
}

func NewProduct(productService product.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		productId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			web.Error(c, 400, "invalid ID")
			return
		}

		ctx := context.Background()
		product, err := h.productService.Get(ctx, productId)

		if err != nil {
			web.Error(c, 404, err.Error())
			return
		}

		web.Success(c, 200, product)
	}
}
