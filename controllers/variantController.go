package controllers

import (
	"final-assignment/initializers"
	"final-assignment/models"
	"final-assignment/models/request"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateVariant(ctx *gin.Context) {
	db := initializers.DB

	var variantReq request.VariantRequest
	if err := ctx.ShouldBind(&variantReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var getProducts models.Product
	if err := db.Model(&getProducts).Where("uuid = ?", variantReq.ProductId).First(&getProducts).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Product Not Found",
		})
		return
	}
	Product := models.Variant{
		VariantName: variantReq.VariantName,
		Quantity:    variantReq.Quantity,
		ProductID:   getProducts.ID,
	}
	err := db.Debug().Create(&Product).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": Product,
	})

}

func GetAllVariant(ctx *gin.Context) {
	db := initializers.DB

	var variants []models.Variant
	if err := db.Preload("Products").Find(&variants).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    variants,
	})
}

func GetOneVariant(ctx *gin.Context) {
	variantId, _ := ctx.Params.Get("variantId")
	var variant models.Variant
	if err := initializers.DB.Preload("Products").Where("uuid = ?", variantId).First(&variant).Error; err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"data":    variant,
	})
}
