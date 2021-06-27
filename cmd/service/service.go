package service

import (
	"github.com/alexdemen/companies/app"
	"net/http"
)

type CompanyService struct {
	CompanyLogic *app.CompanyLogic
}

func (cs *CompanyService) AddBuilding(w http.ResponseWriter, r *http.Request) error {
	cs.CompanyLogic.AddBuilding()
	return nil
}
