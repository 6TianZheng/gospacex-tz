package request

type ProductAdd struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Images      string  `json:"images"`
	Description string  `json:"description"`
	CategoryId  int     `json:"categoryId"`
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
