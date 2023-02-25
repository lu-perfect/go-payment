package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	db "gobank/internal/db/sqlc"
	"gobank/internal/util/validaton"
	"log"
	"net/http"
)

// TODO: move to configuration

const (
	DBDriver      = "postgres"
	DBSource      = "postgresql://root:secret@localhost:5432/gobank?sslmode=disable"
	ServerAddress = "localhost:8080"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer() *Server {
	router := gin.New()

	conn := connectToDB()
	store := db.NewSQLStore(conn)

	s := &Server{
		store:  store,
		router: router,
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
		err := v.RegisterValidation("currency", validaton.CurrencyValidator)
		if err != nil {
			log.Fatal("cannot register currency validator: ", err)
		}
	}
}

func (s *Server) setupRouter() {
	api := s.router.Group("/api")
	{
		accounts := api.Group("/accounts")
		{
			accounts.GET("/:id", s.handleGetAccountById)
			accounts.POST("/", s.handleCreateAccount)
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

func handleInternalServerError(ctx *gin.Context, err error) {
	handleError(ctx, err, http.StatusInternalServerError)
}
