package request

import "mime/multipart"

type CreateProductRequest struct {
	Name  string                `form:"name" binding:"required"`
	Image *multipart.FileHeader `form:"file"`
}

type UpdateProductRequest struct {
	Name  string                `form:"name"`
	Image *multipart.FileHeader `form:"file"`
}
