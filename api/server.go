package api

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	db "github.com/hiiamanop/simple_bank/db/sqlc"
	"github.com/hiiamanop/simple_bank/token"
	"github.com/hiiamanop/simple_bank/util"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	config     util.Config
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {

	fmt.Printf("Server initialized with token duration: %v\n", config.AccessTokenDuration)

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		router:     gin.Default(),
	}

	// Add CORS middleware with more permissive settings
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	// Important: Enable CORS preflight requests by allowing wildcard origins
	// This is necessary for handling requests from different origins during development
	corsConfig.AllowWildcard = true
	corsConfig.MaxAge = 12 * time.Hour

	server.router.Use(cors.New(corsConfig))

	// setup routing
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := server.router

	// Group routes under /api/v1
	v1 := router.Group("/api/v1")

	// Public routes
	users := v1.Group("/users")
	{
		users.POST("", server.createUser)      // signup
		users.POST("/login", server.loginUser) // login
	}

	// Routes that require authentication
	authRoutes := v1.Group("")
	authRoutes.Use(authMiddleware(server.tokenMaker))
	{
		// Protected user routes
		authUsers := authRoutes.Group("/users")
		{
			authUsers.GET("/:username", server.getUser)
			// authUsers.GET("", server.listUsers)
			// authUsers.PUT("/:id", server.updateUser)
			// authUsers.DELETE("/:id", server.deleteUser)
		}

		// Account routes - all protected
		accounts := authRoutes.Group("/accounts")
		{
			accounts.POST("", server.createAccount)
			accounts.GET("/:id", server.getAccount)
			accounts.GET("", server.listAccounts)
			// accounts.PUT("/:id", server.updateAccount)
			// accounts.DELETE("/:id", server.deleteAccount)
		}

		// Entry routes - all protected
		entries := authRoutes.Group("/entries")
		{
			entries.POST("", server.createEntry)
			entries.GET("/:id", server.getEntry)
			entries.GET("", server.listEntries)
			entries.PUT("/:id", server.updateEntry)
			entries.DELETE("/:id", server.deleteEntry)
		}

		// Transfer routes - all protected
		transfers := authRoutes.Group("/transfers")
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
