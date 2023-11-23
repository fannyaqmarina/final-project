package request

type CreateVariantRequest struct {
	VariantName string `form:"variant_name" validate:"required"`
	Quantity    int    `form:"quantity" validate:"required"`
	ProductId   string `form:"product_id" validate:"required"`
}

type UpdateVariantRequest struct {
	VariantName string `form:"variant_name"`
	Quantity    int    `form:"quantity"`
	ProductId   string `form:"product_id"`
}
