package services

import (
	"escrolla-api/config"
	"escrolla-api/db"
	"escrolla-api/models"
)

var apiKey = "sk_test_04f212ad5ac8b0674ac970eb31fa2cd9473b3105"
var escrowFeePercentage = 2 // Your escrow fee percentage
type transactionsService struct {
	Config           *config.Config
	transactionsRepo db.TransactionsRepo
	mail             Mailer
}

func NewTransactionsService(transactionsRepo db.TransactionsRepo, conf *config.Config, mailer Mailer) TransactionsService {
	return &transactionsService{
		Config:           conf,
		transactionsRepo: transactionsRepo,
		mail:             mailer,
	}
}

type TransactionsService interface {
	CreateOrder(paymentRequest models.Order, user models.User) (*models.Order, error)
}

func (t transactionsService) CreateOrder(orderRequest models.Order, user models.User) (*models.Order, error) {
	// Store order information in the database
	order := models.Order{
		BuyerPhone:    orderRequest.BuyerPhone,
		SellerPhone:   orderRequest.SellerPhone,
		BuyerEmail:    orderRequest.BuyerEmail,
		SellerEmail:   orderRequest.SellerEmail,
		Amount:        orderRequest.Amount,
		Description:   orderRequest.Description,
		DeliveryDays:  orderRequest.DeliveryDays,
		UserType:      orderRequest.UserType,
		OrderStatus:   models.Pending,
		PaymentStatus: models.Pending,
		EscrowFee:     orderRequest.EscrowFee,
		UserID:        user.ID,
	}
	orderCreated, err := t.transactionsRepo.CreateOrder(&order)
	if err != nil {
		return nil, err
	}
	return orderCreated, nil
}

/*
	// Calculate escrow fee
	escrowFee := (escrowFeePercentage * paymentRequest.Amount) / 100

	// Create a Paystack client
	httpClient := &http.Client{}
	client := paystack.NewClient(apiKey, httpClient)

	// Create a customer
	customer := &paystack.Customer{
		Email:     paymentRequest.CustomerEmail,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
	}
	customer, err := client.Customer.Create(customer)
	if err != nil {
		gologger.Println(err, "Failed to create customer")
		return nil, fmt.Errorf("an error occured")
	}

	// Create a recipient (freelancer)
	recipient := &paystack.TransferRecipient{
		Type:          "Nuban",
		Name:          paymentRequest.FreelancerName,
		Description:   "Freelancer description",
		AccountNumber: paymentRequest.FreelancerAccountNumber,
		BankCode:      "FREELANCER_BANK_CODE",
		Currency:      "NGN",
		// Add more recipient details as needed
	}
	recipient, err = client.Transfer.CreateRecipient(recipient)
	if err != nil {
		gologger.Println(err, "Failed to create recipient")
		return nil, fmt.Errorf("an error occured")
	}

	// Initiate the transfer to the freelancer
	transferRequest := &paystack.TransferRequest{
		Source:    "balance",
		Reason:    "Payment for service",
		Amount:    float32(paymentRequest.Amount - escrowFee), // Deduct escrow fee
		Recipient: recipient.RecipientCode,
	}
	transactionInfo, err := client.Transfer.Initiate(transferRequest)
	if err != nil {
		gologger.Println(err, "Failed to initiate transfer")
		return nil, fmt.Errorf("an error occured")
	}

*/
