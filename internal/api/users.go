package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	db "gobank/internal/db/sqlc"
	"time"
)

type getUserByIdRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) handleGetUserById(ctx *gin.Context) {
	var req getUserByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		handleBadRequest(ctx, err)
		return
	}

	user, err := s.store.GetUser(ctx, req.ID)

	if err == sql.ErrNoRows {
		handleNotFound(ctx, err)
		return
	}

	if err != nil {
		handleInternalServerError(ctx, err)
		return
	}

	handleSuccess(ctx, user)
}

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (s *Server) handleCreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleBadRequest(ctx, err)
		return
	}

	// TODO: hash password

	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		// TODO: handle db uniq error
		handleInternalServerError(ctx, err)
		return
	}

	res := newUserResponse(user)
	handleCreated(ctx, res)
}

type userResponse struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}
