package request

type VariantRequest struct {
	VariantName string `form:"variant_name" binding:"required"`
	Quantity    int    `form:"quantity"`
	ProductId   string `form:"product_id" binding:"required"`
}
