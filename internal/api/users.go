package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"gobank/internal/auth"
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

	authPayload := getAuthPayload(ctx)
	if user.ID != authPayload.UserID {
		err := errors.New("access denied")
		handleForbidden(ctx, err)
		return
	}

	res := newUserResponse(user)
	handleSuccess(ctx, res)
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

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		handleInternalServerError(ctx, err)
		return
	}

	// TODO check role

	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	})

	if isDBUniqueError(err) {
		handleForbidden(ctx, err)
		return
	}

	if err != nil {
		handleInternalServerError(ctx, err)
		return
	}

	res := newUserResponse(user)
	handleCreated(ctx, res)
}

type userResponse struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}
