package routes

import (
	"github.com/alexdemen/companies/app"
	"github.com/alexdemen/companies/app/stores"
	"github.com/alexdemen/companies/cmd/service"
	"github.com/alexdemen/companies/platform/router"
	"net/http"
)

func API(store stores.Store) http.Handler {
	cs := service.CompanyService{
		CompanyLogic: app.NewCompanyLogic(store),
	}

	handler := router.NewRouter()
	handler.Handle(http.MethodPost, "/v1/building", cs.AddBuilding)
	return handler
}
