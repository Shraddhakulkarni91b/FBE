package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetPointsHandler_InvalidURLFormat(t *testing.T) {
	req, err := http.NewRequest("GET", "/invalid-url", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	GetPointsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetPointsHandler_ReceiptNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/receipts/non-existent-id/points", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	GetPointsHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetPointsHandler_ValidReceiptID(t *testing.T) {
	// Set up a test receipt in storage
	id := "test-id"
	points := 10
	Store.SaveReceipt(id, Receipt{}, points)

	req, err := http.NewRequest("GET", "/receipts/test-id/points", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	GetPointsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]int
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response["points"] != points {
		t.Errorf("expected points %d, got %d", points, response["points"])
	}
}

func TestProcessReceiptHandler(t *testing.T) {
	tests := []struct {
		name       string
		receipt    string
		statusCode int
	}{
		{
			name: "Valid receipt JSON payload",
			receipt: `{
				 "retailer": "Target",
    			 "purchaseDate": "2022-01-01",
    			 "purchaseTime": "13:01",
    			 "items": [
        				{
				 			"shortDescription": "Mountain Dew 12PK", "price": "6.49"
						}
    				],
    			 "total": "35.35"
				 }
			}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "Invalid receipt JSON payload",
			receipt:    `Invalid JSON`,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Receipt with missing fields",
			receipt: `{
				"retailer": "Test Retailer"
			}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Receipt with invalid fields",
			receipt: `{
				"retailer": "Test Retailer",
				"total": " invalid total",
				"items": [
					{
						"shortDescription": "Item 1",
						"price": " invalid price"
					}
				]
			}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Successful receipt processing",
			receipt: `{
				 "retailer": "Target",
    			 "purchaseDate": "2022-01-01",
    			 "purchaseTime": "13:01",
    			 "items": [
        				{"shortDescription": "Mountain Dew 12PK", "price": "6.49"}
    					],
    			 "total": "35.35"}
				]
			}`,
			statusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/receipts/process", strings.NewReader(test.receipt))
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			ProcessReceiptHandler(w, req)

			if w.Code != test.statusCode {
				t.Errorf("expected status code %d, got %d", test.statusCode, w.Code)
			}
		})
	}
}
