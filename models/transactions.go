package models

const Completed = "completed"
const Pending = "pending"

type Transaction struct {
	Model
	Amount                  int
	CustomerEmail           string
	FreelancerName          string
	FreelancerAccountNumber string
	RecipientCode           string
	Integration             int
	Source                  string
	AmountFromPaystack      float32
	Currency                string
	Reason                  string
	TransferCode            string
	Status                  string
	TransferredAt           string
	TitanCode               string
}

type Order struct {
	Model
	UserID        string `json:"user_id"`
	BuyerPhone    string `json:"buyer_phone"`
	SellerPhone   string `json:"seller_phone"`
	BuyerEmail    string `json:"buyer_email"`
	SellerEmail   string `json:"seller_email"`
	Amount        int64  `json:"amount"`
	Description   string `json:"description"`
	DeliveryDays  int    `json:"delivery_days"`
	UserType      string `json:"user_type"` //buyer or seller
	OrderStatus   string `json:"order_status"`
	PaymentStatus string `json:"payment_status"`
	EscrowFee     int64
}

// if usertype is buyer [form must provide seller phone number and email] and viceversa
