package internal

import "sync"

type InMemoryStorage struct {
	receipts map[string]Receipt
	points   map[string]int
	mu       sync.Mutex
}

var Store = InMemoryStorage{
	receipts: make(map[string]Receipt),
	points:   make(map[string]int),
}

// SaveReceipt stores the given receipt and associated points in memory.
// It locks the mutex to ensure thread-safe access to the storage maps.
// Parameters:
//   - id: A unique identifier for the receipt.
//   - receipt: The Receipt object to be stored.
//   - points: The points calculated and associated with the receipt.

func (s *InMemoryStorage) SaveReceipt(id string, receipt Receipt, points int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.receipts[id] = receipt
	s.points[id] = points
}

// GetPoints retrieves the points associated with a given receipt ID.
// It locks the mutex to ensure thread-safe access to the storage maps.
// Parameters:
//   - id: A unique identifier for the receipt.
// Returns:
//   - points: The points calculated and associated with the receipt, if found.
//   - exists: A boolean indicating whether a receipt with the given ID was found.
func (s *InMemoryStorage) GetPoints(id string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	points, exists := s.points[id]
	return points, exists
}
