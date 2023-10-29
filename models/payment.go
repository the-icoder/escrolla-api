package models

type PaymentRequest struct {
	Amount                  int    `json:"amount"`
	CustomerEmail           string `json:"customer_email"`
	FreelancerName          string `json:"freelancer_name"`
	FreelancerAccountNumber string `json:"freelancer_account_number"`
}

type PayStackPaymentRequest struct {
	Amount      int64  `json:"amount"`
	Email       string `json:"email"`
	Reference   string `json:"reference"`
	CallbackUrl string `json:"callback_url"`
}
