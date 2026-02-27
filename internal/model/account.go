package model

// Account represents a bank account.
// JSON tags match the Kotlin Account data class for API contract fidelity.
type Account struct {
	ID         int    `json:"id"`
	Balance    int    `json:"balance"`
	Number     string `json:"number"`
	CustomerID int    `json:"customerId"`
}
