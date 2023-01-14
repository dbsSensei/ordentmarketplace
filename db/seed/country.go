package seed

import (
	"context"
	"fmt"
	db "github.com/dbssensei/ordentmarketplace/db/sqlc"
	"github.com/dbssensei/ordentmarketplace/json"
)

func (s Seed) Country(store db.Store) {
	countries := json.Countries()
	totalCountry, err := store.CountCountry(context.Background())
	if err != nil {
		fmt.Errorf("error when get totalCountry")
	}

	if totalCountry == 0 {
		for _, country := range countries {
			countryParams := db.CreateCountryParams{
				Code: country.Code,
				Name: country.Name,
			}
			store.CreateCountry(context.Background(), countryParams)
		}
		totalCountry, err = store.CountCountry(context.Background())
		if err != nil {
			fmt.Errorf("error when get totalCountry")
		} else if int(totalCountry) < len(countries) {
			fmt.Println("country seeder incomplete")
		} else {
			fmt.Println("country seeder successful")
		}
	}
}
