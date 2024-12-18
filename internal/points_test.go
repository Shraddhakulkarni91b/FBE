package internal

import (
	"testing"
)

func TestCalculatePoints_RetailerName(t *testing.T) {
	receipt := Receipt{
		Retailer: "Test Retailer",
		Total:    "10.00",
	}
	points := CalculatePoints(receipt)
	if points < 10 {
		t.Errorf("expected points to be at least 10, got %d", points)
	}
}

func TestCalculatePoints_RoundDollarAmount(t *testing.T) {
	receipt := Receipt{
		Retailer: "Test Retailer",
		Total:    "10.00",
	}
	points := CalculatePoints(receipt)
	if points < 50 {
		t.Errorf("expected points to be at least 50, got %d", points)
	}
}

func TestCalculatePoints_MultipleOf25Cents(t *testing.T) {
	receipt := Receipt{
		Retailer: "Test Retailer",
		Total:    "10.25",
	}
	points := CalculatePoints(receipt)
	if points < 25 {
		t.Errorf("expected points to be at least 25, got %d", points)
	}
}

func TestCalculatePoints_Items(t *testing.T) {
	receipt := Receipt{
		Retailer: "Test Retailer",
		Total:    "10.00",
		Items: []Item{
			{ShortDescription: "Item 1", Price: "5.00"},
			{ShortDescription: "Item 2", Price: "5.00"},
		},
	}
	points := CalculatePoints(receipt)
	if points < 5 {
		t.Errorf("expected points to be at least 5, got %d", points)
	}
}

func TestCalculatePoints_ItemDescriptionMultipleOf3(t *testing.T) {
	receipt := Receipt{
		Retailer: "Test Retailer",
		Total:    "10.00",
		Items: []Item{
			{ShortDescription: "abc", Price: "5.00"},
		},
	}
	points := CalculatePoints(receipt)
	if points < 1 {
		t.Errorf("expected points to be at least 1, got %d", points)
	}
}

func TestCalculatePoints_TotalGreaterThan10(t *testing.T) {
	receipt := Receipt{
		Retailer: "Test Retailer",
		Total:    "11.00",
	}
	points := CalculatePoints(receipt)
	if points < 5 {
		t.Errorf("expected points to be at least 5, got %d", points)
	}
}

func TestCalculatePoints_OddPurchaseDate(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		Total:        "10.00",
		PurchaseDate: "2022-01-01",
	}
	points := CalculatePoints(receipt)
	if points < 6 {
		t.Errorf("expected points to be at least 6, got %d", points)
	}
}

func TestCalculatePoints_PurchaseTimeBetween2And4(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Test Retailer",
		Total:        "10.00",
		PurchaseTime: "14:00",
	}
	points := CalculatePoints(receipt)
	if points < 10 {
		t.Errorf("expected points to be at least 10, got %d", points)
	}
}
