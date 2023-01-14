package gapi

import (
	"fmt"

	db "github.com/dbssensei/ordentmarketplace/db/sqlc"
	"github.com/dbssensei/ordentmarketplace/pb"
	"github.com/dbssensei/ordentmarketplace/token"
	"github.com/dbssensei/ordentmarketplace/util"
)

// Server serves gRPC requests for our marketplace service.
type Server struct {
	pb.UnimplementedOrdentMarketplaceServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server.
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

	return server, nil
}
