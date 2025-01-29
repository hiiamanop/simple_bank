package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestListAccounts(t *testing.T) {
	// Create 10 random accounts
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestAddAccountBalance(t *testing.T) {
	account1 := createRandomAccount(t)

	// Test cases to verify different balance update scenarios
	testCases := []struct {
		name        string
		amount      int64
		shouldError bool
	}{
		{
			name:        "add positive amount",
			amount:      int64(100),
			shouldError: false,
		},
		{
			name:        "subtract valid amount",
			amount:      int64(-50),
			shouldError: false,
		},
		{
			name:        "add zero",
			amount:      int64(0),
			shouldError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the current balance before update
			currentBalance := account1.Balance

			arg := AddAccountBalanceParams{
				ID:     account1.ID,
				Amount: tc.amount,
			}

			updatedAccount, err := testQueries.AddAccountBalance(context.Background(), arg)

			if tc.shouldError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, updatedAccount)
			require.Equal(t, account1.ID, updatedAccount.ID)
			require.Equal(t, account1.Owner, updatedAccount.Owner)
			require.Equal(t, account1.Currency, updatedAccount.Currency)
			require.Equal(t, currentBalance+tc.amount, updatedAccount.Balance)
			require.Equal(t, account1.CreatedAt, updatedAccount.CreatedAt)

			// Update reference account for next test case
			account1 = updatedAccount
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	// Verify account is deleted
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

// Helper function to create a random account
func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    randomString(6),
		Balance:  randomInt(0, 1000),
		Currency: randomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

// Helper functions for generating random data
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[randomInt(0, int64(len(letters)-1))]
	}
	return string(b)
}

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomCurrency() string {
	currencies := []string{"USD", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
