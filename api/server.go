package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	db "github.com/hiiamanop/simple_bank/db/sqlc"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{
		store:  store,
		router: gin.Default(),
	}

	// Add CORS middleware with more permissive settings
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	// Important: Enable CORS preflight requests
	config.AllowWildcard = true
	config.MaxAge = 12 * time.Hour

	server.router.Use(cors.New(config))

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
			accounts.POST("", server.createAccount)
			accounts.GET("/:id", server.getAccount)
			accounts.GET("", server.listAccounts)
			accounts.PUT("/:id", server.updateAccount)
			accounts.DELETE("/:id", server.deleteAccount)
		}

		// Entry routes
		entries := v1.Group("/entries")
		{
			entries.POST("", server.createEntry)
			entries.GET("/:id", server.getEntry)
			entries.GET("", server.listEntries)
			entries.PUT("/:id", server.updateEntry)
			entries.DELETE("/:id", server.deleteEntry)
		}

		// Transfer routes
		transfers := v1.Group("/transfers")
		{
			transfers.POST("", server.createTransfer)
			transfers.GET("/:id", server.getTransfer)
			transfers.GET("", server.listTransfers)
			transfers.PUT("/:id", server.updateTransfer)
			transfers.DELETE("/:id", server.deleteTransfer)
		}
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
