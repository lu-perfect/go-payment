package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	db "gobank/internal/db/sqlc"
	"log"
	"net/http"
)

// TODO: move to configuration

const (
	DBDriver = "postgres"
	DBSource = "postgresql://root:secret@localhost:5432/gobank?sslmode=disable"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer() *Server {
	router := gin.New()

	conn, err := sql.Open(DBDriver, DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewSQLStore(conn)

	s := &Server{
		store:  store,
		router: router,
	}

	api := router.Group("/api")
	{
		accounts := api.Group("/accounts")
		{
			accounts.GET("/:id", s.handleGetAccountById)
			accounts.POST("/", s.handleCreateAccount)
		}
	}

	return s
}

func (s *Server) Run(address string) error {
	return s.router.Run(address)
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
