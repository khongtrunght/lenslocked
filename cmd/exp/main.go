package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/khongtrunght/lenslocked/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load env: %s", err)
	}

	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("failed to convert port to int: %s", err)
	}
	es, err := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatalf("failed to create email service: %s", err)
	}

	if err := es.ForgotPassword("khongtrunght@gmail.com", "https://google.com"); err != nil {
		log.Fatalf("failed to send forgot password email: %s", err)
	}

	fmt.Println("forgot password email sent")
}
