package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/zhifaq/simple_bank/sqlc"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	if account, err := s.store.CreateAccount(ctx, arg); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	} else {
		ctx.JSON(http.StatusOK, account)
	}
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(ctx *gin.Context) {
	var req GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, account)
}


type ListAccountsRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}
func (s *Server) listAccounts(ctx *gin.Context) {
	var req ListAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}

	accounts, err := s.store.GetAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}