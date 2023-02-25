package api

import (
	"github.com/gin-gonic/gin"
	"gobank/internal/auth"
	db "gobank/internal/db/sqlc"
	"time"
)

const (
	AccessTokenDuration  = time.Minute * 15
	RefreshTokenDuration = time.Hour * 360
)

type signUpRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (s *Server) handleSignUp(ctx *gin.Context) {
	var req signUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleBadRequest(ctx, err)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		handleInternalServerError(ctx, err)
		return
	}

	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	})
	if err != nil {
		// TODO: handle db uniq error
		handleInternalServerError(ctx, err)
		return
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(user.ID, user.Username, AccessTokenDuration)
	if err != nil {
		handleInternalServerError(ctx, err)
		return
	}

	refreshToken, refreshTokenPayload, err := s.tokenMaker.CreateToken(user.ID, user.Username, RefreshTokenDuration)
	if err != nil {
		handleInternalServerError(ctx, err)
		return
	}

	res := newAuthResponse(user, accessToken, accessTokenPayload.ExpiredAt, refreshToken, refreshTokenPayload.ExpiredAt)
	handleCreated(ctx, res)
}

type authResponse struct {
	// SessionID uuid.UUID `json:"session_id"` TODO

	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`

	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`

	User userResponse `json:"user"`
}

func newAuthResponse(user db.User, accessToken string, accessTokenExpiredAt time.Time, refreshToken string, refreshTokenExpiredAt time.Time) authResponse {
	return authResponse{
		// SessionID:             session.ID, TODO
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiredAt,
		User:                  newUserResponse(user),
	}
}
