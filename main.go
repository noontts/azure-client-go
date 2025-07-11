package main

import (
	"azureclient/client/azure"
	http_client "azureclient/client/http"
	"azureclient/config"
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := azure.NewAzureClient(cfg.Azure.StorageAccount, cfg.Azure.EmailEndpoint, cfg.Azure.EmailAccessKey)
	if err != nil {
		log.Fatalf("Failed to create AzureClient: %v", err)
	}

	ctx := context.Background()
	// Example: Upload a blob
	err = client.BlobClient.UploadBlob(ctx, "test-container", "test-blob.txt", []byte("Hello from AzureClient!"))
	if err != nil {
		fmt.Println("UploadBlob error:", err)
	} else {
		fmt.Println("Blob uploaded successfully.")
	}

	// Example: Download a blob
	data, err := client.BlobClient.DownloadBlob(ctx, "test-container", "test-blob.txt")
	if err != nil {
		fmt.Println("DownloadBlob error:", err)
	} else {
		fmt.Println("Downloaded blob:", string(data))
	}

	// Example: Send an email
	emailReq := azure.SendEmailRequest{
		Sender:    "sender@yourdomain.com",
		Recipient: "recipient@other.com",
		Subject:   "Test Subject",
		PlainText: "Hello from AzureClient Email!",
		HTML:      "<b>Hello from AzureClient Email!</b>",
	}
	err = client.SendEmailClient.SendEmail(ctx, emailReq)
	if err != nil {
		fmt.Println("SendEmail error:", err)
	} else {
		fmt.Println("Email sent successfully.")
	}

	// Example: PaymentClient usage
	paymentClient := http_client.NewPaymentClient(cfg.Client.PaymentBaseURL)
	if err := paymentClient.DoSomething(ctx); err != nil {
		fmt.Println("PaymentClient error:", err)
	} else {
		fmt.Println("PaymentClient call succeeded.")
	}

	// Example: MemberClient usage
	memberClient := http_client.NewMemberClient(cfg.Client.MemberBaseURL)
	if err := memberClient.DoSomething(ctx); err != nil {
		fmt.Println("MemberClient error:", err)
	} else {
		fmt.Println("MemberClient call succeeded.")
	}
}
