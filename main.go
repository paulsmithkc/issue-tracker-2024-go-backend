package main

import (
	"log"
	"os"

	// Blank-import the function package so the init() runs
	// _ "example.com/gcf/v2"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/joho/godotenv"
)

func main() {
  // Load environment variables from .env file
  godotenv.Load(".env")

	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// By default, listen on all interfaces. If testing locally, run with
	// LOCAL_ONLY=true to avoid triggering firewall warnings and
	// exposing the server outside of your own machine.
	// hostname := ""
	// if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
	// 	hostname = "127.0.0.1"
	// }
  hostname := "127.0.0.1"

  //fmt.Fprintln(log, "OK")
  log.Printf("Listening on http://%s:%s\n", hostname, port)

	if err := funcframework.StartHostPort(hostname, port); err != nil {
		log.Fatalf("funcframework.StartHostPort: %v\n", err)
	}
}
