package repository

import (
	"testing"

	"github.com/sample-kotlin-ktor-microservices/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddAccount(t *testing.T) {
	repo := NewInMemoryAccountRepository()

	account := repo.Add(model.Account{Balance: 1000, Number: "1234", CustomerID: 1})

	assert.Equal(t, 1, account.ID)
	assert.Equal(t, 1000, account.Balance)
	assert.Equal(t, "1234", account.Number)
	assert.Equal(t, 1, account.CustomerID)
}

func TestAddMultipleAccounts(t *testing.T) {
	repo := NewInMemoryAccountRepository()

	a1 := repo.Add(model.Account{Balance: 100, Number: "111", CustomerID: 1})
	a2 := repo.Add(model.Account{Balance: 200, Number: "222", CustomerID: 2})

	assert.Equal(t, 1, a1.ID)
	assert.Equal(t, 2, a2.ID)
}

func TestGetAll(t *testing.T) {
	repo := NewInMemoryAccountRepository()
	repo.Add(model.Account{Balance: 100, Number: "111", CustomerID: 1})
	repo.Add(model.Account{Balance: 200, Number: "222", CustomerID: 2})

	accounts := repo.GetAll()

	assert.Len(t, accounts, 2)
}

func TestGetAllEmpty(t *testing.T) {
	repo := NewInMemoryAccountRepository()

	accounts := repo.GetAll()

	assert.Empty(t, accounts)
	assert.NotNil(t, accounts)
}

func TestGetByID(t *testing.T) {
	repo := NewInMemoryAccountRepository()
	repo.Add(model.Account{Balance: 1000, Number: "1234", CustomerID: 1})

	account, err := repo.GetByID(1)

	require.NoError(t, err)
	assert.Equal(t, 1, account.ID)
	assert.Equal(t, 1000, account.Balance)
}

func TestGetByIDNotFound(t *testing.T) {
	repo := NewInMemoryAccountRepository()

	_, err := repo.GetByID(999)

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestGetByCustomerID(t *testing.T) {
	repo := NewInMemoryAccountRepository()
	repo.Add(model.Account{Balance: 100, Number: "111", CustomerID: 1})
	repo.Add(model.Account{Balance: 200, Number: "222", CustomerID: 2})
	repo.Add(model.Account{Balance: 300, Number: "333", CustomerID: 1})

	accounts := repo.GetByCustomerID(1)

	assert.Len(t, accounts, 2)
	for _, a := range accounts {
		assert.Equal(t, 1, a.CustomerID)
	}
}

func TestGetByCustomerIDNoMatch(t *testing.T) {
	repo := NewInMemoryAccountRepository()
	repo.Add(model.Account{Balance: 100, Number: "111", CustomerID: 1})

	accounts := repo.GetByCustomerID(999)

	assert.Empty(t, accounts)
	assert.NotNil(t, accounts)
}
