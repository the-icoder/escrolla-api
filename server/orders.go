package server

import (
	"bytes"
	"encoding/json"
	"escrolla-api/errors"
	"escrolla-api/models"
	"escrolla-api/server/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var escrowFeePercentage int64 = 2                                    // Your escrow fee percentage
const secretKey = "sk_test_04f212ad5ac8b0674ac970eb31fa2cd9473b3105" // Replace with your Paystack secret key
// const esc =
//func (s *Server) Paystack() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		var paymentRequest models.PaymentRequest
//		if err := c.ShouldBindJSON(&paymentRequest); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//		user := models.User{}
//		transactions, err := s.TransactionsService.CreateOrder(paymentRequest, user)
//		if err != nil {
//			response.JSON(c, "", http.StatusBadRequest, nil, err)
//			return
//		}
//		response.JSON(c, "Payment received and held in escrow.", http.StatusOK, transactions, nil)
//	}
//}

func (s *Server) CreateOrderAndPay() gin.HandlerFunc {
	return func(c *gin.Context) {
		//_, user, err := GetValuesFromContext(c)
		//if err != nil {
		//	response.JSON(c, "", err.Status, nil, err)
		//	return
		//}
		var orderRequest models.Order
		if err := c.ShouldBindJSON(&orderRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		escrowFee := (escrowFeePercentage * orderRequest.Amount) / 100
		paymentRequest := models.PayStackPaymentRequest{
			Amount: (orderRequest.Amount + escrowFee) * 100,
			Email:  orderRequest.BuyerEmail,
		}
		user := models.User{}
		orderRequest.EscrowFee = escrowFee
		order, err := s.TransactionsService.CreateOrder(orderRequest, user)
		if err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		paymentRequest.Reference = order.ID
		resp, err := sendPaystackRequest("POST", &paymentRequest)
		if err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		log.Println("checking reference: ", order.ID, paymentRequest.Reference)
		//transactions, err := s.TransactionsService.CreateOrder(paymentRequest, user)
		//if err != nil {
		//	response.JSON(c, "", http.StatusBadRequest, nil, err)
		//	return
		//}
		response.JSON(c, "initializing payment.", http.StatusOK, resp, nil)

	}
}

func (s *Server) VerifyPayment(c *gin.Context) {
	reference := c.DefaultQuery("reference", "")
	if reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No reference supplied"})
		return
	}

	// Verify the Paystack transaction
	verifyURL := "https://api.paystack.co/transaction/verify/" + reference

	resp, err := sendPaystackVerifyRequest("GET", verifyURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resp.Status {
		err = s.TransactionsRepo.UpdateOrderStatus(reference, "paid")
		if err != nil {
			response.JSON(c, "kindly await the status of your payment", http.StatusInternalServerError, resp, nil)
			return
		}
		response.JSON(c, "Congratulations, your Escrow deposit successfully locked", http.StatusOK, resp, nil)
	} else {
		// Handle the case when the transaction was not successful
		response.JSON(c, resp.Message, http.StatusBadRequest, resp, nil)
	}
}
func sendPaystackVerifyRequest(method, url string, payload map[string]interface{}) (*PaystackResponse, error) {
	client := &http.Client{}

	var reqBody []byte
	if payload != nil {
		reqBody, _ = json.Marshal(payload)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Paystack request failed with status code: %d", resp.StatusCode)
	}

	var response PaystackResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
func sendPaystackRequest(method string, payload *models.PayStackPaymentRequest) (*PaystackResponse, error) {
	client := &http.Client{}
	payload.CallbackUrl = "https://warm-frangipane-c63d2f.netlify.app/"
	var reqBody []byte
	if payload != nil {
		reqBody, _ = json.Marshal(payload)
	}
	initURL := "https://api.paystack.co/transaction/initialize"
	req, err := http.NewRequest(method, initURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Paystack request failed with status code: %d", resp.StatusCode)
	}

	var response PaystackResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

type PaystackResponse struct {
	Status bool `json:"status"`
	Data   struct {
		AuthorizationURL string `json:"authorization_url"`
		Status           string `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
}

func (s *Server) GetOrderByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		//_, user, err := GetValuesFromContext(c)
		//if err != nil {
		//	err.Respond(c)
		//	return
		//}
		//
		orders, err := s.TransactionsRepo.GetOrderByUserID("1")
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("internal server error", http.StatusInternalServerError))
			return
		}
		response.JSON(c, "retrieved orders successfully", http.StatusOK, gin.H{"orders": orders}, nil)
	}
}
