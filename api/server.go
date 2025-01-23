package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/hiiamanop/simple_bank/db/sqlc"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{
		store:  store,
		router: gin.Default(),
	}

	// setup routing
	server.setupRouter()

	return server
}

func (server *Server) setupRouter() {
	// Add routes to the router
	router := server.router

	// Group routes under /api/v1
	v1 := router.Group("/api/v1")
	{
		// Account routes
		accounts := v1.Group("/accounts")
		{
			accounts.POST("/", server.createAccount)
			accounts.GET("/:id", server.getAccount)
			accounts.GET("/", server.listAccounts)
			accounts.PUT("/:id", server.updateAccount)
			accounts.DELETE("/:id", server.deleteAccount)
		}

		// Add more route groups here as needed
		// Example:
		// transfers := v1.Group("/transfers")
		// entries := v1.Group("/entries")
	}
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse is a helper function to return error responses
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
