package main

import (
	"escrolla-api/config"
	"escrolla-api/db"
	"escrolla-api/server"
	"escrolla-api/services"
	"log"
	"net/http"
	"time"
)

func main() {
	http.DefaultClient.Timeout = time.Second * 10
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	gormDB := db.GetDB(conf)
	authRepo := db.NewAuthRepo(gormDB)
	transactionsRepo := db.NewTransactions(gormDB)
	mail := services.NewMailService(conf)
	authService := services.NewAuthService(authRepo, conf, mail)
	transactionsService := services.NewTransactionsService(transactionsRepo, conf, mail)

	s := &server.Server{
		Config:              conf,
		AuthRepository:      authRepo,
		AuthService:         authService,
		TransactionsRepo:    transactionsRepo,
		TransactionsService: transactionsService,
	}
	s.Start()
}
