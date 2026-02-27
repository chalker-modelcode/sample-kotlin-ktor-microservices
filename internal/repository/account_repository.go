package repository

import (
	"errors"
	"sync"

	"github.com/sample-kotlin-ktor-microservices/internal/model"
)

// ErrNotFound is returned when a requested entity does not exist.
var ErrNotFound = errors.New("not found")

// AccountRepository defines the interface for account storage operations.
type AccountRepository interface {
	// Add stores a new account and returns it with an auto-assigned ID.
	Add(account model.Account) model.Account
	// GetAll returns all stored accounts.
	GetAll() []model.Account
	// GetByID returns the account with the given ID, or ErrNotFound.
	GetByID(id int) (model.Account, error)
	// GetByCustomerID returns all accounts belonging to the given customer.
	GetByCustomerID(customerID int) []model.Account
}

// InMemoryAccountRepository is a thread-safe in-memory implementation of AccountRepository.
type InMemoryAccountRepository struct {
	mu       sync.RWMutex
	accounts []model.Account
	nextID   int
}

// NewInMemoryAccountRepository creates a new empty in-memory account repository.
func NewInMemoryAccountRepository() *InMemoryAccountRepository {
	return &InMemoryAccountRepository{
		accounts: make([]model.Account, 0),
		nextID:   1,
	}
}

// Add stores a new account, assigns an auto-incremented ID, and returns the stored account.
func (r *InMemoryAccountRepository) Add(account model.Account) model.Account {
	r.mu.Lock()
	defer r.mu.Unlock()

	account.ID = r.nextID
	r.nextID++
	r.accounts = append(r.accounts, account)
	return account
}

// GetAll returns a copy of all stored accounts.
func (r *InMemoryAccountRepository) GetAll() []model.Account {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.Account, len(r.accounts))
	copy(result, r.accounts)
	return result
}

// GetByID returns the account with the given ID, or ErrNotFound if it does not exist.
func (r *InMemoryAccountRepository) GetByID(id int) (model.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, a := range r.accounts {
		if a.ID == id {
			return a, nil
		}
	}
	return model.Account{}, ErrNotFound
}

// GetByCustomerID returns all accounts belonging to the given customer.
func (r *InMemoryAccountRepository) GetByCustomerID(customerID int) []model.Account {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.Account
	for _, a := range r.accounts {
		if a.CustomerID == customerID {
			result = append(result, a)
		}
	}
	if result == nil {
		result = make([]model.Account, 0)
	}
	return result
}
