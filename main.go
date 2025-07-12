package main

import (
	"azureclient/client/azure"
	http_client "azureclient/client/http"
	"azureclient/config"
	"azureclient/internal/controller"
	"azureclient/internal/repository"
	"azureclient/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"

	"azureclient/internal/otel"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormotel "gorm.io/plugin/opentelemetry"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// --- OpenTelemetry Initialization ---
	shutdown, err := otel.InitTracer()
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// GORM MySQL connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Add OpenTelemetry plugin for GORM
	if err := db.Use(gormotel.NewPlugin()); err != nil {
		log.Fatalf("Failed to register GORM OpenTelemetry plugin: %v", err)
	}

	// DI: Repositories struct, Service, Controller
	repos := repository.Repositories{
		Member: repository.NewMemberRepository(db),
		// Add more repositories here as needed
	}
	memberService := service.NewMemberService(repos)
	memberController := controller.NewMemberController(memberService)

	// Set up Gorilla Mux router
	r := mux.NewRouter()
	memberController.RegisterRoutes(r)

	// Start HTTP server
	go func() {
		fmt.Println("HTTP server started on :8080")
		if err := http.ListenAndServe(":8080", r); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// --- Existing Azure/Client logic below ---
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
		Sender:    "DoNotReply@9cb54951-6182-4a23-8dee-328302efdcf2.azurecomm.net",
		Recipient: "noonthitisan@gmail.com",
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
