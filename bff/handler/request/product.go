package request

type ProductAdd struct {
	Name        string  `form:"name"  binding:"required"`
	Price       float64 `form:"price"  binding:"required"`
	Images      string  `form:"images"  binding:"required"`
	Description string  `form:"description"  binding:"required"`
	CategoryId  int     `form:"categoryId"  binding:"required"`
}

type OrderItem struct {
	ProductId int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

type CreateOrderRequest struct {
	UserId    int64       `json:"user_id"`
	PayType   int64       `json:"pay_type"`
	AddressId int64       `json:"address_id"`
	List      []OrderItem `json:"list"`
}
