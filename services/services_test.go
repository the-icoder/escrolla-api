package services

import (
	"escrolla-api/config"
	"fmt"
	"log"
	"os"
	"testing"
)

var testConfig *config.Config

func TestMain(m *testing.M) {
	fmt.Println("Starting server tests")
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	testConfig = c
	fmt.Println(testConfig)
	exitCode := m.Run()
	os.Exit(exitCode)
}
