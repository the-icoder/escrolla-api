package server

import (
	"escrolla-api/models"
	"escrolla-api/server/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) Paystack() gin.HandlerFunc {
	return func(c *gin.Context) {
		var paymentRequest models.PaymentRequest
		if err := c.ShouldBindJSON(&paymentRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		transactions, err := s.TransactionsService.CreateTransactions(paymentRequest)
		if err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		response.JSON(c, "Payment received and held in escrow.", http.StatusOK, transactions, nil)
	}
}
