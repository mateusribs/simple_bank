package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/mateusribs/simple_bank/db/sqlc"
)

// serves HTTP requests for banking services
type Server struct {
	store db.Store
	router *gin.Engine
}


// create a new server instance
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/accounts/update", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	server.router = router

	return server
}

// runs HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}