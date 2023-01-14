package api

import (
	"fmt"

	db "github.com/dbssensei/ordentmarketplace/db/sqlc"
	"github.com/dbssensei/ordentmarketplace/token"
	"github.com/dbssensei/ordentmarketplace/util"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our marketplace service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	gin.SetMode(server.config.GinMode)
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)
	router.GET("/categories", server.getCategories)
	router.GET("/categories/:id", server.getCategory)
	router.GET("/products", server.getProducts)
	router.GET("/products/:id", server.getProduct)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.PATCH("/users", server.updateUser)
	authRoutes.POST("/merchants", server.createMerchant)
	authRoutes.PATCH("/merchants", server.updateMerchant)
	authRoutes.POST("/categories", server.createCategory)
	authRoutes.DELETE("/categories/:id", server.deleteCategory)
	authRoutes.POST("/products", server.createProduct)
	authRoutes.PATCH("/products/:id", server.updateProduct)
	authRoutes.DELETE("/products/:id", server.deleteProduct)
	authRoutes.GET("/orders", server.getOrders)
	authRoutes.GET("/orders/:id", server.getOrder)
	authRoutes.POST("/orders", server.createOrder)
	authRoutes.PATCH("/orders", server.updateOrderStatus)

	// merchant_id : 1
	// product_id : 2
	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
