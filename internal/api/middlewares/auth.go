package middlewares

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gobank/internal/auth/token"
	"net/http"
	"strings"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

var ErrAuthHeaderNotProvided = errors.New("authorization header is not provided")
var ErrInvalidAuthHeaderFormat = errors.New("invalid authorization header format")

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			handleAbortWithUnauthorized(ctx, ErrAuthHeaderNotProvided)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			handleAbortWithUnauthorized(ctx, ErrInvalidAuthHeaderFormat)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			handleAbortWithUnauthorized(ctx, fmt.Errorf("unsupported authorization type %s", authorizationType))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			handleAbortWithUnauthorized(ctx, err)
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}

func handleAbortWithUnauthorized(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": err.Error(),
	})
}
