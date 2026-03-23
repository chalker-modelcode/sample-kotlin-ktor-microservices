package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/sample-kotlin-ktor-microservices/internal/model"
	"github.com/sample-kotlin-ktor-microservices/internal/repository"
)

// AccountHandler handles HTTP requests for account operations.
type AccountHandler struct {
	repo repository.AccountRepository
}

// NewAccountHandler creates a new AccountHandler with the given repository.
func NewAccountHandler(repo repository.AccountRepository) *AccountHandler {
	return &AccountHandler{repo: repo}
}

// CreateAccount handles POST /accounts.
// It decodes the request body as an Account, stores it, and returns the created account.
// Returns HTTP 200 to match the original Kotlin service behavior.
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account model.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created := h.repo.Add(account)
	writeJSON(w, http.StatusOK, created)
}

// ListAccounts handles GET /accounts.
// It returns all accounts as a JSON array.
func (h *AccountHandler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	accounts := h.repo.GetAll()
	writeJSON(w, http.StatusOK, accounts)
}

// GetAccountByID handles GET /accounts/{id}.
// It returns the account with the given ID, or 404 if not found.
func (h *AccountHandler) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	account, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "account not found")
		return
	}

	writeJSON(w, http.StatusOK, account)
}

// GetAccountsByCustomerID handles GET /accounts/customer/{customerId}.
// It returns all accounts belonging to the given customer as a JSON array.
func (h *AccountHandler) GetAccountsByCustomerID(w http.ResponseWriter, r *http.Request) {
	customerIDStr := chi.URLParam(r, "customerId")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	accounts := h.repo.GetByCustomerID(customerID)
	writeJSON(w, http.StatusOK, accounts)
}
