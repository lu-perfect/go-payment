package server

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func New() *Server {
	router := gin.New()

	return &Server{
		router: router,
	}
}

func (s *Server) Run(address string) error {
	return s.router.Run(address)
}
