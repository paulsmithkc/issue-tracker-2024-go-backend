package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
  log.Println("Init healthCheck")
  functions.HTTP("healthCheck", healthCheck)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
  log.Println("Request healthCheck")
  fmt.Fprintln(w, "OK")
}
