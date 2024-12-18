package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// ProcessReceiptHandler handles HTTP requests for processing receipts.
// It decodes the JSON payload to a Receipt, calculates points based on the receipt,
// generates a unique ID, saves the receipt and associated points to storage, and
// responds with the generated ID as a JSON object.

func ProcessReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		log.Printf("Error decoding receipt payload: %v", err)
		http.Error(w, "Invalid receipt format. Please verify input.", http.StatusBadRequest)
		return
	}
	// Validate receipt payload
	valid, validationErrors := validateReceipt(receipt)
	if !valid {
		log.Printf("Receipt validation failed: %v", validationErrors)
		response := ValidationResponse{
			Message: "Validation Failed",
			Errors:  validationErrors,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	id := uuid.New().String()
	log.Printf("Generated receipt ID: %s", id)
	points := CalculatePoints(receipt)
	log.Printf("Calculated points for receipt ID %s: %d", id, points)
	Store.SaveReceipt(id, receipt, points)

	log.Printf("Saved receipt with ID %s to storage", id)
	response := map[string]string{"id": id}
	json.NewEncoder(w).Encode(response)
}

// GetPointsHandler handles HTTP requests for getting points associated with a receipt.
// It expects the receipt ID to be passed as part of the URL path, as /receipts/{id}/points.
// It fetches the points for the given ID from storage and responds with the points as JSON.
// If the receipt is not found, it returns a 404. If the URL format is invalid, it returns a 400.
func GetPointsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s", r.URL.Path)
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 || pathParts[1] != "receipts" || pathParts[3] != "points" {
		log.Printf("Invalid URL format. Use /receipts/{id}/points.")
		http.Error(w, "Invalid URL format. Use /receipts/{id}/points.", http.StatusBadRequest)
		return
	}

	id := pathParts[2]
	// Fetch the points for the given ID from storage
	log.Printf("Fetching points for receipt ID: %s", id)
	points, exists := Store.GetPoints(id)
	if !exists {
		log.Printf("Receipt with ID %s not found.", id)
		http.Error(w, "Receipt not found.", http.StatusNotFound)
		return
	}

	// Respond with the points as JSON
	log.Printf("Responding with points for receipt ID %s: %d", id, points)
	response := map[string]int{"points": points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func validateReceipt(receipt Receipt) (bool, []string) {
	var errors []string

	// Validate retailer
	if receipt.Retailer == "" {
		errors = append(errors, "Retailer is required.")
	} else {
		retailerPattern := regexp.MustCompile(`^[\w\s\-&]+$`)
		if !retailerPattern.MatchString(receipt.Retailer) {
			errors = append(errors, "Retailer contains invalid characters.")
		}
	}

	// Validate purchaseDate
	if receipt.PurchaseDate == "" {
		errors = append(errors, "PurchaseDate is required.")
	} else {
		datePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		if !datePattern.MatchString(receipt.PurchaseDate) {
			errors = append(errors, "PurchaseDate must be in YYYY-MM-DD format.")
		}
	}

	// Validate purchaseTime
	if receipt.PurchaseTime == "" {
		errors = append(errors, "PurchaseTime is required.")
	} else {
		timePattern := regexp.MustCompile(`^\d{2}:\d{2}$`)
		if !timePattern.MatchString(receipt.PurchaseTime) {
			errors = append(errors, "PurchaseTime must be in HH:MM format.")
		}
	}

	// Validate total
	if receipt.Total == "" {
		errors = append(errors, "Total is required.")
	} else {
		totalPattern := regexp.MustCompile(`^\d+\.\d{2}$`)
		if !totalPattern.MatchString(receipt.Total) {
			errors = append(errors, "Total must be a decimal value with two decimal places.")
		}
	}

	// Validate items
	if len(receipt.Items) == 0 {
		errors = append(errors, "At least one item is required.")
	} else {
		for i, item := range receipt.Items {
			if item.ShortDescription == "" {
				errors = append(errors, "Item "+strconv.Itoa(i+1)+" is missing a ShortDescription.")
			}
			if item.Price == "" {
				errors = append(errors, "Item "+strconv.Itoa(i+1)+" is missing a Price.")
			} else {
				pricePattern := regexp.MustCompile(`^\d+\.\d{2}$`)
				if !pricePattern.MatchString(item.Price) {
					errors = append(errors, "Item "+strconv.Itoa(i+1)+" Price must be a decimal value with two decimal places.")
				}
			}
		}
	}

	return len(errors) == 0, errors
}
