package main

import (
	"log"
	"net/http"
	"receipt-processor/internal"
)

// Main is the entry point for the receipt processor server.
// It sets up two routes: /receipts/process and /receipts/{id}/points.
// The former is used to process receipts, and the latter is used to retrieve points associated with a receipt.
// The server listens on port 8080.
func main() {
	http.HandleFunc("/receipts/process", internal.ProcessReceiptHandler)
	http.HandleFunc("/receipts/", internal.GetPointsHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
