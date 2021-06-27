package main

import (
	"context"
	"github.com/alexdemen/companies/cmd/routes"
	"github.com/alexdemen/companies/platform/postgres"
	"net/http"
)

func main() {
	pgInstance := postgres.NewInstance(context.Background())
	routes.API(pgInstance)

	api := http.Server{
		Addr:    ":8080",
		Handler: routes.API(pgInstance),
	}

	if err := api.ListenAndServe(); err != nil {

	}
}
