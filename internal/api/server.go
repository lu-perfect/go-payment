package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"gobank/internal/api/middlewares"
	"gobank/internal/auth/token"
	db "gobank/internal/db/sqlc"
	"gobank/internal/util"
	"log"
	"net/http"
)

// TODO: move to configuration

const (
	DBDriver          = "postgres"
	DBSource          = "postgresql://root:secret@localhost:5432/gobank?sslmode=disable"
	ServerAddress     = "localhost:8080"
	TokenSymmetricKey = "NiIsInR5cCI6IgRG9lIiwiaWF0IjoxlK" // 32
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer() *Server {
	router := gin.New()

	conn := connectToDB()
	store := db.NewSQLStore(conn)
	tokenMaker := createTokenMaker()

	s := &Server{
		store:      store,
		router:     router,
		tokenMaker: tokenMaker,
	}

	s.registerValidators()
	s.setupRouter()

	return s
}

func (s *Server) Run() error {
	return s.router.Run(ServerAddress)
}

func (s *Server) registerValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", util.CurrencyValidator)
		if err != nil {
			log.Fatal("cannot register currency validator: ", err)
		}
	}
}

func (s *Server) setupRouter() {
	authMiddleware := middlewares.AuthMiddleware(s.tokenMaker)

	api := s.router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/sign-up", s.handleSignUp)
			auth.POST("/sign-in", s.handleSignIn)
		}

		accounts := api.Group("/accounts")
		accounts.Use(authMiddleware)
		{
			accounts.GET("/:id", s.handleGetAccountById)
			accounts.POST("", s.handleCreateAccount)
		}

		users := api.Group("/users")
		users.Use(authMiddleware)
		{
			users.GET("/:id", s.handleGetUserById)
			users.POST("", s.handleCreateUser)
		}
	}
}

func connectToDB() *sql.DB {
	conn, err := sql.Open(DBDriver, DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	return conn
}

func createTokenMaker() token.Maker {
	tokenMaker, err := token.NewPasetoMaker(TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker: %w", err)
	}
	return tokenMaker
}

func getAuthPayload(ctx *gin.Context) *token.Payload {
	return ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
}

func handleSuccess(ctx *gin.Context, obj any) {
	ctx.JSON(http.StatusOK, obj)
}

func handleCreated(ctx *gin.Context, obj any) {
	ctx.JSON(http.StatusCreated, obj)
}

func handleError(ctx *gin.Context, err error, code int) {
	ctx.JSON(code, gin.H{
		"error": err.Error(),
	})
}

func handleBadRequest(ctx *gin.Context, err error) {
	handleError(ctx, err, http.StatusBadRequest)
}

func handleNotFound(ctx *gin.Context, err error) {
	handleError(ctx, err, http.StatusNotFound)
}

func handleForbidden(ctx *gin.Context, err error) {
	handleError(ctx, err, http.StatusForbidden)
}

func handleUnauthorized(ctx *gin.Context, err error) {
	handleError(ctx, err, http.StatusUnauthorized)
}

func handleInternalServerError(ctx *gin.Context, err error) {
	handleError(ctx, err, http.StatusInternalServerError)
}
