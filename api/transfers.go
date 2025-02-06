package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/hiiamanop/simple_bank/db/sqlc"
	"github.com/hiiamanop/simple_bank/token"
)

type createTransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gt=0"`
}

type transferResponse struct {
	Transfer    db.Transfer `json:"transfer"`
	FromAccount db.Account  `json:"from_account"`
	ToAccount   db.Account  `json:"to_account"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get the authenticated user from the context
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Validate from account exists and belongs to the authenticated user
	fromAccount, err := server.store.GetAccount(ctx, req.FromAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Verify account ownership
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Check if from_account has sufficient balance
	if fromAccount.Balance < req.Amount {
		err := fmt.Errorf("account %d has insufficient balance for transfer: %d < %d",
			req.FromAccountID, fromAccount.Balance, req.Amount)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validate to account exists
	toAccount, err := server.store.GetAccount(ctx, req.ToAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Validate currencies match
	if fromAccount.Currency != toAccount.Currency {
		err := fmt.Errorf(
			"currency mismatch: from account [%d] currency %s vs to account [%d] currency %s",
			req.FromAccountID,
			fromAccount.Currency,
			req.ToAccountID,
			toAccount.Currency,
		)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := transferResponse{
		Transfer:    result.Transfer,
		FromAccount: result.FromAccount,
		ToAccount:   result.ToAccount,
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) getTransfer(ctx *gin.Context) {
	var req struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get the authenticated user
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Get the transfer
	transfer, err := server.store.GetTransfers(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get the from account to verify ownership
	fromAccount, err := server.store.GetAccount(ctx, transfer.FromAccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Verify the user owns either the from or to account
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("transfer doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

type listTransfersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listTransfers(ctx *gin.Context) {
	var req listTransfersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListTransfersParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	transfers, err := server.store.ListTransfers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfers)
}

type updateTransferRequest struct {
	Amount int64 `json:"amount" binding:"required,gt=0"`
}

func (server *Server) updateTransfer(ctx *gin.Context) {
	var reqURI struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqBody updateTransferRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get the authenticated user
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Get the transfer to verify ownership
	transfer, err := server.store.GetTransfers(ctx, reqURI.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get the from account to verify ownership
	fromAccount, err := server.store.GetAccount(ctx, transfer.FromAccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Verify ownership
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("transfer doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateTransferParams{
		ID:     reqURI.ID,
		Amount: reqBody.Amount,
	}

	transfer, err = server.store.UpdateTransfer(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

func (server *Server) deleteTransfer(ctx *gin.Context) {
	var req struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get the authenticated user
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Get the transfer to verify ownership
	transfer, err := server.store.GetTransfers(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get the from account to verify ownership
	fromAccount, err := server.store.GetAccount(ctx, transfer.FromAccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Verify ownership
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("transfer doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteTransfers(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Transfer deleted successfully"})
}
