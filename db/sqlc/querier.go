// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
)

type Querier interface {
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error)
	// Create a new account
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	// Create a new entries
	CreateEntries(ctx context.Context, arg CreateEntriesParams) (Entry, error)
	// Create a new transfers
	CreateTransfers(ctx context.Context, arg CreateTransfersParams) (Transfer, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	// Delete an account
	DeleteAccount(ctx context.Context, id int64) error
	// Delete an entries
	DeleteEntries(ctx context.Context, id int64) error
	// Delete a transfers
	DeleteTransfers(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, username string) error
	// Get an account by id
	GetAccount(ctx context.Context, id int64) (Account, error)
	// Get an entries by id
	GetEntries(ctx context.Context, id int64) (Entry, error)
	// Get a transfers by id
	GetTransfers(ctx context.Context, id int64) (Transfer, error)
	GetUser(ctx context.Context, username string) (User, error)
	// List all accounts
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	// List all entries
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
	// List all transfers
	ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateEntries(ctx context.Context, arg UpdateEntriesParams) (Entry, error)
	UpdateTransfer(ctx context.Context, arg UpdateTransferParams) (Transfer, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
