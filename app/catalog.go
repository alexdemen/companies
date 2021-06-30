package app

import (
	"github.com/alexdemen/companies/app/models"
	"github.com/alexdemen/companies/app/stores"
)

type Catalog struct {
	store stores.Store
}

func NewCatalog(store stores.Store) *Catalog {
	return &Catalog{
		store: store,
	}
}

func (cl *Catalog) AddBuilding(b models.Building) (id int64, err error) {
	ex, err := cl.store.GetExecutor()
	if err != nil {
		return
	}
	defer func() {
		ex.Close(err)
	}()

	b.Id, err = ex.AddBuilding(b)
	if err != nil {
		return
	}

	if b.Organizations != nil && len(*b.Organizations) > 0 {
		for _, org := range *b.Organizations {
			org.Building = &models.Building{Id: b.Id}
			if _, err = ex.AddOrganization(org); err != nil {
				return
			}
		}
	}

	return
}

func (cl *Catalog) GetOrganizationByIds(ids []int64) (res []models.Organization, err error) {
	ex, err := cl.store.GetExecutor()
	if err != nil {
		return
	}
	defer func() {
		ex.Close(err)
	}()

	err = ex.GetOrganizationsByIds(&res, ids)

	return
}

func (cl *Catalog) GetOrganizationInBuilding(buildingId int64) (res []models.Organization, err error) {
	ex, err := cl.store.GetExecutor()
	if err != nil {
		return
	}
	defer func() {
		ex.Close(err)
	}()

	err = ex.GetOrganizationInBuilding(&res, buildingId)

	return
}

func (cl *Catalog) GetOrganizationByCategory(categoryId int64) (res []models.Organization, err error) {
	ex, err := cl.store.GetExecutor()
	if err != nil {
		return
	}
	defer func() {
		ex.Close(err)
	}()

	err = ex.GetOrganizationByCategory(&res, categoryId)

	return
}

func (cl *Catalog) GetCategoryList() (res []models.Category, err error) {
	ex, err := cl.store.GetExecutor()
	if err != nil {
		return nil, err
	}
	defer func() {
		ex.Close(err)
	}()

	err = ex.GetCategoryList(&res)

	return
}
