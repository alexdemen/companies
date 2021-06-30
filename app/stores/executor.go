package stores

import "github.com/alexdemen/companies/app/models"

type Executor interface {
	Close(err error)
	AddBuilding(b models.Building) (int64, error)
	AddOrganization(org models.Organization) (int64, error)
	GetOrganizationsByIds(orgs *[]models.Organization, ids []int64) error
	GetOrganizationInBuilding(orgs *[]models.Organization, buildingId int64) error
	GetOrganizationByCategory(orgs *[]models.Organization, categoryId int64) error
	GetCategoryList(data *[]models.Category) error
}
