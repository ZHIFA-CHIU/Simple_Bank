package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/zhifaq/simple_bank/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store:  store,
		router: gin.Default(),
	}

	// Define routes here
	server.router.POST("/accounts", server.createAccount)
	server.router.GET("/accounts/:id", server.getAccount)
	server.router.GET("/accounts", server.listAccounts)

	return server
}

// Start runs the HTTP server on the specified address.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
