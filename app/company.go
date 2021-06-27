package app

import "github.com/alexdemen/companies/app/stores"

type CompanyLogic struct {
	store stores.Store
}

func NewCompanyLogic(store stores.Store) *CompanyLogic {
	return &CompanyLogic{
		store: store,
	}
}

func (cl *CompanyLogic) AddBuilding() {
	ex := cl.store.GetExecutor()
	defer ex.Close(nil)
}
