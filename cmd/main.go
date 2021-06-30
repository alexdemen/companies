package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/alexdemen/companies/cmd/routes"
	"github.com/alexdemen/companies/platform/postgres"
)

func main() {
	pgInstance, err := postgres.NewInstance(context.Background(), os.Getenv("POSTGRES_CONNECTION_URL")) //"postgres://postgres:12345678@127.0.0.1:5433/companies")
	if err != nil {
		log.Fatal(err)
	}
	defer pgInstance.Close()

	api := http.Server{
		Addr:    ":8080",
		Handler: routes.API(pgInstance),
	}

	if err := api.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
