package model

type Order struct {
	OrderID string      `json:"id" gorm:"primaryKey;type:varchar(64)"`
	UserID  int         `json:"user_id"`
	Orders  []OrderItem `gorm:"foreignKey:OrderID;references:OrderID"`
	Address string      `json:"address"`
	Total   float64     `json:"total"`
}

type OrderItem struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	OrderID   string `json:"order_id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type OrderMessage struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
