package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	db "gobank/internal/db/sqlc"
)

type getAccountByIdRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) handleGetAccountById(ctx *gin.Context) {
	var req getAccountByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		handleBadRequest(ctx, err)
		return
	}

	account, err := s.store.GetAccount(ctx, req.ID)

	if err == sql.ErrNoRows {
		handleNotFound(ctx, err)
		return
	}

	if err != nil {
		handleInternalServerError(ctx, err)
		return
	}

	handleSuccess(ctx, account)
}

type createAccountRequest struct {
	// TODO: validate available currencies
	Currency string `json:"currency" binding:"required"`
	Owner    string `json:"owner" binding:"required"`
}

func (s *Server) handleCreateAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleBadRequest(ctx, err)
		return
	}

	account, err := s.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
	})

	if err != nil {
		// TODO: handle unique db error
		handleInternalServerError(ctx, err)
	}

	handleCreated(ctx, account)
}
