package models

type Order struct {
	OrderID string  `json:"order_id"`
	Product string  `json:"product"`
	Qty     int     `json:"qty"`
	Price   float64 `json:"price"`
}

func MockOrders() []Order {
	orders := []Order{
		{
			OrderID: "487e4d63-ff55-43af-987f-1bfbedc4a57f",
			Product: "iPhone 7",
			Qty:     1,
			Price:   150000,
		},
		{
			OrderID: "e59fdcf8-b256-40f1-be63-d31834f87150",
			Product: "iPhone 8",
			Qty:     3,
			Price:   250000,
		},
	}

	return orders
}
