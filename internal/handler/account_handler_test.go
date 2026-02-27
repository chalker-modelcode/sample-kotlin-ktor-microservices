package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sample-kotlin-ktor-microservices/internal/model"
	"github.com/sample-kotlin-ktor-microservices/internal/repository"
)

// setupAccountRouter creates a chi router wired with an AccountHandler
// backed by a fresh in-memory repository, suitable for testing.
func setupAccountRouter() (*chi.Mux, *repository.InMemoryAccountRepository) {
	repo := repository.NewInMemoryAccountRepository()
	h := NewAccountHandler(repo)

	r := chi.NewRouter()
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", h.CreateAccount)
		r.Get("/", h.ListAccounts)
		r.Get("/{id}", h.GetAccountByID)
		r.Get("/customer/{customerId}", h.GetAccountsByCustomerID)
	})
	return r, repo
}

// createTestAccount is a helper that POSTs an account and returns the response.
func createTestAccount(t *testing.T, router *chi.Mux, balance int, number string, customerID int) *httptest.ResponseRecorder {
	t.Helper()
	body := map[string]interface{}{
		"balance":    balance,
		"number":     number,
		"customerId": customerID,
	}
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// TestCreateAccount mirrors Kotlin's testAdd - POST /accounts creates an account.
func TestCreateAccount(t *testing.T) {
	router, _ := setupAccountRouter()

	w := createTestAccount(t, router, 1000, "1234", 1)

	assert.Equal(t, http.StatusOK, w.Code)

	var account model.Account
	err := json.Unmarshal(w.Body.Bytes(), &account)
	require.NoError(t, err)
	assert.Equal(t, 1, account.ID)
	assert.Equal(t, 1000, account.Balance)
	assert.Equal(t, "1234", account.Number)
	assert.Equal(t, 1, account.CustomerID)
}

// TestListAccounts mirrors Kotlin's testFindAll - GET /accounts returns all accounts.
func TestListAccounts(t *testing.T) {
	router, _ := setupAccountRouter()

	// Create two accounts
	createTestAccount(t, router, 1000, "1234", 1)
	createTestAccount(t, router, 2000, "5678", 2)

	// List all
	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var accounts []model.Account
	err := json.Unmarshal(w.Body.Bytes(), &accounts)
	require.NoError(t, err)
	assert.Len(t, accounts, 2)
}

// TestListAccountsEmpty verifies GET /accounts returns an empty array when no accounts exist.
func TestListAccountsEmpty(t *testing.T) {
	router, _ := setupAccountRouter()

	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var accounts []model.Account
	err := json.Unmarshal(w.Body.Bytes(), &accounts)
	require.NoError(t, err)
	assert.Empty(t, accounts)
}

// TestGetAccountByID mirrors Kotlin's testFindById - GET /accounts/{id} returns account.
func TestGetAccountByID(t *testing.T) {
	router, _ := setupAccountRouter()

	// Create an account first
	createTestAccount(t, router, 1000, "1234", 1)

	// Get by ID
	req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var account model.Account
	err := json.Unmarshal(w.Body.Bytes(), &account)
	require.NoError(t, err)
	assert.Equal(t, 1, account.ID)
	assert.Equal(t, 1000, account.Balance)
	assert.Equal(t, "1234", account.Number)
	assert.Equal(t, 1, account.CustomerID)
}

// TestGetAccountByIDNotFound verifies GET /accounts/{id} returns 404 for missing accounts.
func TestGetAccountByIDNotFound(t *testing.T) {
	router, _ := setupAccountRouter()

	req := httptest.NewRequest(http.MethodGet, "/accounts/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var errResp errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "account not found", errResp.Error)
}

// TestGetAccountByIDInvalidID verifies GET /accounts/{id} returns 400 for non-numeric ID.
func TestGetAccountByIDInvalidID(t *testing.T) {
	router, _ := setupAccountRouter()

	req := httptest.NewRequest(http.MethodGet, "/accounts/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errResp errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "invalid account id", errResp.Error)
}

// TestGetAccountsByCustomerID verifies GET /accounts/customer/{customerId} filters correctly.
func TestGetAccountsByCustomerID(t *testing.T) {
	router, _ := setupAccountRouter()

	// Create accounts for two customers
	createTestAccount(t, router, 100, "111", 1)
	createTestAccount(t, router, 200, "222", 2)
	createTestAccount(t, router, 300, "333", 1)

	// Get accounts for customer 1
	req := httptest.NewRequest(http.MethodGet, "/accounts/customer/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var accounts []model.Account
	err := json.Unmarshal(w.Body.Bytes(), &accounts)
	require.NoError(t, err)
	assert.Len(t, accounts, 2)
	for _, a := range accounts {
		assert.Equal(t, 1, a.CustomerID)
	}
}

// TestGetAccountsByCustomerIDNoMatch verifies empty array for unknown customer.
func TestGetAccountsByCustomerIDNoMatch(t *testing.T) {
	router, _ := setupAccountRouter()

	req := httptest.NewRequest(http.MethodGet, "/accounts/customer/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var accounts []model.Account
	err := json.Unmarshal(w.Body.Bytes(), &accounts)
	require.NoError(t, err)
	assert.Empty(t, accounts)
}

// TestCreateAccountInvalidJSON verifies POST /accounts returns 400 for invalid JSON.
func TestCreateAccountInvalidJSON(t *testing.T) {
	router, _ := setupAccountRouter()

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errResp errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "invalid request body", errResp.Error)
}

// TestCreateAccountAutoIncrementIDs verifies IDs are sequentially assigned.
func TestCreateAccountAutoIncrementIDs(t *testing.T) {
	router, _ := setupAccountRouter()

	w1 := createTestAccount(t, router, 100, "111", 1)
	w2 := createTestAccount(t, router, 200, "222", 2)
	w3 := createTestAccount(t, router, 300, "333", 3)

	var a1, a2, a3 model.Account
	require.NoError(t, json.Unmarshal(w1.Body.Bytes(), &a1))
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &a2))
	require.NoError(t, json.Unmarshal(w3.Body.Bytes(), &a3))

	assert.Equal(t, 1, a1.ID)
	assert.Equal(t, 2, a2.ID)
	assert.Equal(t, 3, a3.ID)
}
