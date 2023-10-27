package models

type PaymentRequest struct {
	Amount                  int    `json:"amount"`
	CustomerEmail           string `json:"customer_email"`
	FreelancerName          string `json:"freelancer_name"`
	FreelancerAccountNumber string `json:"freelancer_account_number"`
}
