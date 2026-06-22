package handler

// ---------------------------------------------------------
// Product Types
// ---------------------------------------------------------

type ProductReq struct {
	Name         string  `json:"name"`
	Image        string  `json:"image"`
	Category     string  `json:"category"`
	Description  *string `json:"description"`
	Rating       int     `json:"rating"`
	NumReviews   int     `json:"num_reviews"`
	Price        float64 `json:"price"`
	CountInStock int     `json:"count_in_stock"`
}

type ProductRes struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Image        string  `json:"image"`
	Category     string  `json:"category"`
	Description  *string `json:"description"`
	Rating       int     `json:"rating"`
	NumReviews   int     `json:"num_reviews"`
	Price        float64 `json:"price"`
	CountInStock int     `json:"count_in_stock"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// ---------------------------------------------------------
// Order Types
// ---------------------------------------------------------

type OrderItemReq struct {
	ProductID int    `json:"product_id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	Image     string `json:"image"`
	Price     int    `json:"price"`
}

type OrderReq struct {
	PaymentMethod string         `json:"payment_method"`
	TaxPrice      float64        `json:"tax_price"`
	ShippingPrice float64        `json:"shipping_price"`
	TotalPrice    float64        `json:"total_price"`
	Items         []OrderItemReq `json:"items"`
}

type OrderItemRes struct {
	ID        int    `json:"id"`
	OrderID   int    `json:"order_id"`
	ProductID int    `json:"product_id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	Image     string `json:"image"`
	Price     int    `json:"price"`
}

type OrderRes struct {
	ID            int            `json:"id"`
	PaymentMethod string         `json:"payment_method"`
	TaxPrice      float64        `json:"tax_price"`
	ShippingPrice float64        `json:"shipping_price"`
	TotalPrice    float64        `json:"total_price"`
	CreatedAt     string         `json:"created_at"`
	UpdatedAt     string         `json:"updated_at"`
	Items         []OrderItemRes `json:"items"`
}

// ---------------------------------------------------------
// General Types
// ---------------------------------------------------------

type DefaultRes struct {
	Message string `json:"message"`
}

type IDRes struct {
	ID int `json:"id"`
}
