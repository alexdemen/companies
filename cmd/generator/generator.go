package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/alexdemen/companies/app"
	"github.com/alexdemen/companies/app/models"
	"github.com/alexdemen/companies/platform/postgres"

	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	orgCount, buildingCount := os.Getenv("ORG_COUNT"), os.Getenv("BUILDING_COUNT")
	oCount, err := strconv.Atoi(orgCount)
	if err != nil {
		fmt.Println(err)
		return
	}
	bCount, err := strconv.Atoi(buildingCount)
	if err != nil {
		fmt.Println(err)
		return
	}

	pgInstance, err := postgres.NewInstance(context.Background(), os.Getenv("POSTGRES_CONNECTION_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pgInstance.Close()

	catalog := app.NewCatalog(pgInstance)
	categoryList, err := catalog.GetCategoryList()
	if err != nil {
		fmt.Println(err)
		return
	}

	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)

	buildings := make([]models.Building, 0, bCount)
	for idx := 0; idx < bCount; idx++ {
		slice := make([]models.Organization, 0, 2)
		buildings = append(buildings, models.Building{
			Address:       faker.Address().Address,
			Latitude:      faker.Latitude(),
			Longitude:     faker.Longitude(),
			Organizations: &slice,
		})
	}

	var org models.Organization
	for idx := 0; idx < oCount; idx++ {
		org = models.Organization{
			Name:         fmt.Sprintf("%s %v", faker.Company(), idx),
			Categories:   []int64{categoryList[faker.Rand.Intn(len(categoryList))].Id},
			PhoneNumbers: []string{faker.Phone()},
		}

		bIdx := faker.Rand.Intn(len(buildings))
		*buildings[bIdx].Organizations = append(*buildings[bIdx].Organizations, org)
	}

	for _, b := range buildings {
		if len(*b.Organizations) > 0 {
			b.Id, err = catalog.AddBuilding(b)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("added: %d organizations", len(*b.Organizations))
		}
	}
}
