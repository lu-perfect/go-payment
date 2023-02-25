package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		UserID:       user.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		handleInternalServerError(ctx, err)
		return
	}

	res := newAuthResponse(session.ID, user, accessToken, accessTokenPayload.ExpiredAt, refreshToken, refreshTokenPayload.ExpiredAt)
	handleCreated(ctx, res)
}

type signInRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

func (s *Server) handleSignIn(ctx *gin.Context) {
	var req signInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleBadRequest(ctx, err)
		return
	}

	user, err := s.store.GetUserByUsername(ctx, req.Username)
	if err == sql.ErrNoRows {
		handleNotFound(ctx, err)
		return
	}
	if err != nil {
		handleInternalServerError(ctx, err)
		return
	}

	err = auth.CheckPassword(req.Password, user.Password)
	if err != nil {
		handleUnauthorized(ctx, err)
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

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		UserID:       user.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		handleInternalServerError(ctx, err)
		return
	}

	res := newAuthResponse(session.ID, user, accessToken, accessTokenPayload.ExpiredAt, refreshToken, refreshTokenPayload.ExpiredAt)
	handleCreated(ctx, res)
}

type refreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type refreshAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (s *Server) handleRefreshAccessToken(ctx *gin.Context) {
	var req refreshAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleBadRequest(ctx, err)
		return
	}

	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		handleUnauthorized(ctx, err)
		return
	}

	session, err := s.store.GetSession(ctx, refreshPayload.ID)
	if err == sql.ErrNoRows {
		handleNotFound(ctx, err)
		return
	}
	if err != nil {
		handleInternalServerError(ctx, err)
		return
	}

	if session.IsRevoked {
		err := fmt.Errorf("blocked session")
		handleUnauthorized(ctx, err)
		return
	}

	if session.UserID != refreshPayload.UserID || session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		handleUnauthorized(ctx, err)
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched session refreshPayload")
		handleUnauthorized(ctx, err)
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		handleUnauthorized(ctx, err)
		return
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(refreshPayload.UserID, refreshPayload.Username, AccessTokenDuration)
	if err != nil {
		handleInternalServerError(ctx, err)
		return
	}

	res := refreshAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
	}
	handleSuccess(ctx, res)
}

type authResponse struct {
	SessionID uuid.UUID `json:"session_id"`

	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`

	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`

	User userResponse `json:"user"`
}

func newAuthResponse(sessionID uuid.UUID, user db.User, accessToken string, accessTokenExpiredAt time.Time, refreshToken string, refreshTokenExpiredAt time.Time) authResponse {
	return authResponse{
		SessionID:             sessionID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiredAt,
		User:                  newUserResponse(user),
	}
}
