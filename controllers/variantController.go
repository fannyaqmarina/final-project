package controllers

import (
	"final-assignment/initializers"
	"final-assignment/models"
	"final-assignment/models/request"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateVariant(ctx *gin.Context) {
	db := initializers.DB

	var variantReq request.CreateVariantRequest
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

	search := ctx.Query("search")
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	var variants []models.Variant
	query := db.Preload("Products")

	if search != "" {
		query = query.Where("variant_name LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	offset, err := strconv.Atoi(limit)
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid limit value",
			"message": err.Error(),
		})
		return
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid page value",
			"message": err.Error(),
		})
		return
	}
	offset = (pageNum - 1) * offset

	if err := query.Offset(offset).Limit(limitInt).Find(&variants).Error; err != nil {
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

func UpdateVariant(ctx *gin.Context) {
	db := initializers.DB

	Variant := models.Variant{}

	variantId := ctx.Param("variantId")

	var variantReq request.UpdateVariantRequest
	if err := ctx.ShouldBind(&variantReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var getVariant models.Variant
	if err := db.Model(&getVariant).Where("uuid = ?", variantId).First(&getVariant).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	Variant.ID = getVariant.ID

	updatedData := models.Variant{
		VariantName: variantReq.VariantName,
		Quantity:    variantReq.Quantity,
	}

	if err := db.Model(&Variant).Where("uuid = ?", variantId).Updates(updatedData).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Success Updated",
	})

}

func DeleteVariant(ctx *gin.Context) {
	db := initializers.DB

	Variant := &models.Variant{}

	variantId := ctx.Param("variantId")

	var getVariant models.Variant
	if err := db.Model(&getVariant).Where("uuid = ?", variantId).First(&getVariant).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	Variant.ID = getVariant.ID
	if err := initializers.DB.Delete(Variant).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete Variant"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Success Deleted Variant",
	})
}
