package seed

import db "github.com/dbssensei/ordentmarketplace/db/sqlc"

type Seed struct {
}

func Execute(store db.Store) {
	var seed Seed
	seed.Country(store)
}
