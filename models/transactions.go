package models

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
