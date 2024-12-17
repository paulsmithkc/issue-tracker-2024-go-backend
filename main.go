package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/joho/godotenv"
)

func main() {
  // Load environment variables from .env file
  godotenv.Load(".env")

  // Use PORT environment variable, or default to 8080
  port := "8080"
  if envPort := os.Getenv("PORT"); envPort != "" {
    port = envPort
  }

  hostname := "127.0.0.1"

  log.Printf("Listening on http://%s:%s\n", hostname, port)

  if err := funcframework.StartHostPort(hostname, port); err != nil {
    log.Fatalf("funcframework.StartHostPort: %v\n", err)
  }
}
