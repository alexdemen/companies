package routes

import (
	"net/http"

	"github.com/alexdemen/companies/app"
	"github.com/alexdemen/companies/app/stores"
	"github.com/alexdemen/companies/cmd/service"
	"github.com/alexdemen/companies/platform/router"
)

func API(store stores.Store) http.Handler {
	cs := service.CompanyService{
		CompanyLogic: app.NewCatalog(store),
	}

	handler := router.NewRouter()
	handler.Handle(http.MethodPost, "/v1/building", cs.AddBuilding)
	handler.Handle(http.MethodPost, "/v1/organization/ids", cs.GetOrganizationByIds)
	handler.Handle(http.MethodGet, "/v1/organization/building/{building_id}", cs.GetOrganizationInBuilding)
	handler.Handle(http.MethodGet, "/v1/organization/category/{category_id}", cs.GetOrganizationByCategory)
	return handler
}
