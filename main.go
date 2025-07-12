package main

import (
	"azureclient/client/azure"
	"azureclient/config"
	"azureclient/internal"
	"azureclient/internal/controller"
	"azureclient/internal/repository"
	"azureclient/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"azureclient/internal/otel"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormotel "gorm.io/plugin/opentelemetry/tracing"
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

	// Azure client for email
	client, err := azure.NewAzureClient(cfg.Azure.StorageAccount, cfg.Azure.EmailEndpoint, cfg.Azure.EmailAccessKey)
	if err != nil {
		log.Fatalf("Failed to create AzureClient: %v", err)
	}
	emailController := controller.NewEmailController(client.SendEmailClient)

	// Set up HTTP routes using SetupRoutes from internal/route.go
	mux := http.NewServeMux()
	internal.SetupRoutes(mux, memberController, emailController)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Graceful shutdown setup
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("HTTP server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for shutdown signal
	sig := <-shutdownCh
	log.Printf("Received signal %s, shutting down HTTP server...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server Shutdown: %v", err)
	}
	log.Println("HTTP server gracefully stopped.")

// Example: Send an email (now handled by API endpoint /send-email)
}
