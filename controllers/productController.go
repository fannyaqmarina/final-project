package controllers

import (
	"final-assignment/helpers"
	"final-assignment/initializers"
	"final-assignment/models"
	"final-assignment/models/request"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

var (
	appJSON = "application/json"
)

func CreateProduct(ctx *gin.Context) {
	db := initializers.DB

	var productReq request.CreateProductRequest
	if err := ctx.ShouldBind(&productReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if productReq.Image == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Image is Required"})
		return
	}

	fileName := helpers.RemoveExtension(productReq.Image.Filename)

	uploadResult, err := helpers.UploadFile(productReq.Image, fileName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userData := ctx.MustGet("userData").(jwt5.MapClaims)

	Product := models.Product{
		Name:     strings.ToLower(productReq.Name),
		ImageUrl: uploadResult,
		AdminID:  uint(userData["id"].(float64)),
	}

	err = db.Debug().Create(&Product).Error
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

func GetAllProduct(ctx *gin.Context) {
	db := initializers.DB

	search := ctx.Query("search")
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	var products []models.Product
	query := db.Preload("Admin").Preload("Variants")

	if search != "" {
		query = query.Where("name LIKE ?", "%"+strings.ToLower(search)+"%")
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

	if err := query.Offset(offset).Limit(limitInt).Find(&products).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    products,
	})
}

func GetOneProduct(ctx *gin.Context) {
	productId, _ := ctx.Params.Get("productId")
	var product models.Product
	if err := initializers.DB.Preload("Variants").Preload("Admin").Where("uuid = ?", productId).First(&product).Error; err != nil {
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
		"data":    product,
	})
}

func UpdateProduct(ctx *gin.Context) {
	db := initializers.DB

	userData := ctx.MustGet("userData").(jwt5.MapClaims)
	Product := models.Product{}

	productId := ctx.Param("productId")
	adminID := uint(userData["id"].(float64))

	var productReq request.UpdateProductRequest
	if err := ctx.ShouldBind(&productReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var getProduct models.Product
	if err := db.Model(&getProduct).Where("uuid = ?", productId).First(&getProduct).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	Product.AdminID = adminID
	Product.ID = getProduct.ID

	updatedData := models.Product{
		Name: productReq.Name,
	}

	if productReq.Image != nil {
		fileName := helpers.RemoveExtension(productReq.Image.Filename)

		uploadResult, err := helpers.UploadFile(productReq.Image, fileName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updatedData.ImageUrl = uploadResult
	}

	if err := db.Model(&Product).Where("uuid = ?", productId).Updates(updatedData).Error; err != nil {
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

func DeleteProduct(ctx *gin.Context) {
	db := initializers.DB

	Product := &models.Product{}

	productId := ctx.Param("productId")

	var getProduct models.Product
	if err := db.Model(&getProduct).Preload("Variants").Where("uuid = ?", productId).First(&getProduct).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	for i := 0; i < len(getProduct.Variants); i++ {
		if err := initializers.DB.Delete(&models.Variant{}, getProduct.Variants[i].ID).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to delete Variant"})
			return
		}
	}

	Product.ID = getProduct.ID
	if err := initializers.DB.Delete(Product).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete Product"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Success Deleted Product",
	})
}
