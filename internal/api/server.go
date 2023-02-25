package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"gobank/internal/api/middlewares"
	"gobank/internal/auth/token"
	db "gobank/internal/db/sqlc"
	"gobank/internal/util"
	"log"
	"net/http"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config) *Server {
	s := &Server{
		config: config,
	}

	s.connectToDB()
	s.addTokenMaker()
	s.registerValidators()
	s.setupRouter()

	return s
}

func (s *Server) Run() error {
	return s.router.Run(s.config.ServerAddress)
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
	router := gin.New()
	authMiddleware := middlewares.AuthMiddleware(s.tokenMaker)

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/sign-up", s.handleSignUp)
			auth.POST("/sign-in", s.handleSignIn)
			auth.POST("/refresh", s.handleRefreshAccessToken)
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

	s.router = router
}

func (s *Server) connectToDB() {
	conn, err := sql.Open(s.config.DBDriver, s.config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	s.store = db.NewSQLStore(conn)
}

func (s *Server) addTokenMaker() {
	tokenMaker, err := token.NewPasetoMaker(s.config.TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker: %w", err)
	}
	s.tokenMaker = tokenMaker
}

func getAuthPayload(ctx *gin.Context) *token.Payload {
	return ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
}

func isDBUniqueError(err error) bool {
	if err == nil {
		return false
	}
	if pqErr, ok := err.(*pq.Error); ok {
		code := pqErr.Code.Name()
		return code == "unique_violation" || code == "foreign_key_violation"
	}
	return false
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
