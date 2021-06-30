package postgres

import (
	"context"
	"strconv"

	"github.com/alexdemen/companies/app/models"
	"github.com/alexdemen/companies/app/stores"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Instance struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func (inst *Instance) GetExecutor() (stores.Executor, error) {
	tx, err := inst.pool.BeginTx(inst.ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	if err != nil {
		return nil, err
	}
	return &Executor{tx: tx, ctx: inst.ctx}, nil
}

func NewInstance(ctx context.Context, connUrl string) (*Instance, error) {
	pool, err := pgxpool.Connect(ctx, connUrl)
	if err != nil {
		return nil, err
	}

	return &Instance{pool: pool, ctx: ctx}, nil
}

func (inst *Instance) Close() {
	inst.pool.Close()
}

type Executor struct {
	tx  pgx.Tx
	ctx context.Context
}

func (e *Executor) Close(err error) {
	if err != nil {
		e.tx.Rollback(e.ctx)
		return
	}

	e.tx.Commit(e.ctx)
}

func (e *Executor) AddOrganization(org models.Organization) (id int64, err error) {
	result, err := e.tx.Query(e.ctx, `
								INSERT INTO organizations
								(name, building_id) VALUES ($1, $2)
								RETURNING id`, org.Name, org.Building.Id)

	if err != nil {
		return 0, err
	}

	if result.Next() {
		err = result.Scan(&id)
	}
	result.Close()
	if err != nil {
		return 0, err
	}

	if len(org.Categories) > 0 {
		err = e.setOrganizationCategories(id, org.Categories)
		if err != nil {
			return
		}
	}

	if len(org.PhoneNumbers) > 0 {
		err = e.setOrganizationPhones(id, org.PhoneNumbers)
		if err != nil {
			return
		}
	}

	return
}

func (e *Executor) AddBuilding(b models.Building) (id int64, err error) {
	result, err := e.tx.Query(e.ctx, `
		INSERT INTO buildings (address, coordinate)	
		VALUES ($1, point($2,$3))
		RETURNING id`, b.Address, b.Latitude, b.Longitude)

	if err != nil {
		return 0, err
	}
	defer result.Close()

	if result.Next() {
		err = result.Scan(&id)
	}

	return
}

func (e *Executor) GetOrganizationsByIds(orgs *[]models.Organization, ids []int64) error {
	var paramIN string
	var params []interface{}
	for idx := range ids {
		if idx > 0 {
			paramIN += ", "
		}

		paramIN += "$" + strconv.Itoa(idx+1)
		params = append(params, ids[idx])
	}
	phones, err := e.getOrganizationPhones(ids)
	if err != nil {
		return err
	}

	categories, err := e.getOrganizationCategories(ids)
	if err != nil {
		return err
	}

	result, err := e.tx.Query(e.ctx, `
								SELECT organizations.id, organizations.name,
									   buildings.id, buildings.address,
									   buildings.coordinate[0] AS latitude, buildings.coordinate[1] AS longitude
								FROM organizations
								INNER JOIN buildings ON buildings.id = organizations.building_id
								WHERE organizations.deleted = false AND organizations.id IN (`+paramIN+`)
	`, params...)
	if err != nil {
		return err
	}
	defer result.Close()

	for result.Next() {
		org := models.Organization{Building: &models.Building{}}
		if err = result.Scan(&org.Id, &org.Name, &org.Building.Id, &org.Building.Address, &org.Building.Latitude, &org.Building.Longitude); err != nil {
			break
		}
		org.PhoneNumbers = phones[org.Id]
		org.Categories = categories[org.Id]
		*orgs = append(*orgs, org)
	}

	return err
}

func (e *Executor) GetOrganizationInBuilding(orgs *[]models.Organization, buildingId int64) error {
	result, err := e.tx.Query(e.ctx, `
								SELECT organizations.id, organizations.name,
									   buildings.id, buildings.address,
									   buildings.coordinate[0] AS latitude, buildings.coordinate[1] AS longitude
								FROM organizations
										 INNER JOIN buildings ON buildings.id = organizations.building_id
								WHERE buildings.id = $1 AND organizations.deleted = false`, buildingId)
	if err != nil {
		return err
	}
	defer result.Close()

	var orgIds []int64
	for result.Next() {
		org := models.Organization{Building: &models.Building{}}
		if err = result.Scan(&org.Id, &org.Name, &org.Building.Id, &org.Building.Address, &org.Building.Latitude, &org.Building.Longitude); err != nil {
			return err
		}
		orgIds = append(orgIds, org.Id)
		*orgs = append(*orgs, org)
	}

	phones, err := e.getOrganizationPhones(orgIds)
	if err != nil {
		return err
	}
	categories, err := e.getOrganizationCategories(orgIds)
	if err != nil {
		return err
	}

	for idx := range *orgs {
		(*orgs)[idx].PhoneNumbers = phones[(*orgs)[idx].Id]
		(*orgs)[idx].Categories = categories[(*orgs)[idx].Id]
	}

	return err
}

func (e *Executor) GetOrganizationByCategory(orgs *[]models.Organization, categoryId int64) error {
	result, err := e.tx.Query(e.ctx, `
							SELECT category_id
							FROM category_hierarchy
							WHERE category_hierarchy.category_up_id = $1`, categoryId)
	if err != nil {
		return err
	}

	ids := []int64{categoryId}
	for result.Next() {
		var id int64
		if err = result.Scan(&id); err != nil {
			break
		}
		ids = append(ids, id)
	}
	result.Close()
	if err != nil {
		return err
	}

	var paramIN string
	var params []interface{}
	for idx := range ids {
		if idx > 0 {
			paramIN += ", "
		}

		paramIN += "$" + strconv.Itoa(idx+1)
		params = append(params, ids[idx])
	}

	result, err = e.tx.Query(e.ctx, `
						SELECT organizations.id, organizations.name,
       					buildings.id, buildings.address,
       					buildings.coordinate[0] AS latitude, buildings.coordinate[1] AS longitude
						FROM  organization_by_category
						INNER JOIN organizations ON organizations.id = organization_by_category.organization_id
						INNER JOIN buildings ON buildings.id = organizations.building_id
						WHERE organization_by_category.category_id IN (`+paramIN+`) AND organizations.deleted = false`, params...)
	if err != nil {
		return err
	}
	defer result.Close()

	var orgIds []int64
	for result.Next() {
		org := models.Organization{Building: &models.Building{}}
		if err = result.Scan(&org.Id, &org.Name, &org.Building.Id, &org.Building.Address, &org.Building.Latitude, &org.Building.Longitude); err != nil {
			return err
		}
		orgIds = append(orgIds, org.Id)
		*orgs = append(*orgs, org)
	}

	phones, err := e.getOrganizationPhones(orgIds)
	if err != nil {
		return err
	}
	categories, err := e.getOrganizationCategories(orgIds)
	if err != nil {
		return err
	}

	for idx := range *orgs {
		(*orgs)[idx].PhoneNumbers = phones[(*orgs)[idx].Id]
		(*orgs)[idx].Categories = categories[(*orgs)[idx].Id]
	}

	return err
}

func (e *Executor) GetCategoryList(data *[]models.Category) error {
	result, err := e.tx.Query(e.ctx, `
			SELECT categories.id, categories.name
			FROM categories
			WHERE categories.deleted = false`)
	if err != nil {
		return err
	}
	defer result.Close()

	var (
		id   int64
		name string
	)
	for result.Next() {
		if err = result.Scan(&id, &name); err != nil {
			break
		}
		*data = append(*data, models.Category{
			Id:   id,
			Name: name,
		})
	}

	return err
}

func (e *Executor) setOrganizationCategories(orgId int64, categories []int64) error {
	for _, category := range categories {
		_, err := e.tx.Exec(e.ctx, `INSERT INTO organization_by_category (organization_id, category_id) VALUES ($1, $2)`,
			orgId, category)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Executor) setOrganizationPhones(orgId int64, phones []string) error {
	for _, phone := range phones {
		_, err := e.tx.Exec(e.ctx, "INSERT INTO organization_phone_numbers (number, organization_id) VALUES ($1, $2)",
			phone, orgId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Executor) getOrganizationPhones(orgIds []int64) (res map[int64][]string, err error) {
	var paramIN string
	var params []interface{}
	for idx := range orgIds {
		if idx > 0 {
			paramIN += ", "
		}

		paramIN += "$" + strconv.Itoa(idx+1)
		params = append(params, orgIds[idx])
	}

	result, err := e.tx.Query(e.ctx, `
					SELECT organization_id, number
					FROM organization_phone_numbers
					WHERE deleted = false AND organization_id IN (`+paramIN+`)`, params...)
	if err != nil {
		return
	}
	defer result.Close()

	res = make(map[int64][]string, len(orgIds))
	var (
		id    int64
		phone string
	)
	for result.Next() {
		err = result.Scan(&id, &phone)
		if err != nil {
			break
		}

		res[id] = append(res[id], phone)
	}

	return
}

func (e *Executor) getOrganizationCategories(orgIds []int64) (res map[int64][]int64, err error) {
	var paramIN string
	var params []interface{}
	for idx := range orgIds {
		if idx > 0 {
			paramIN += ", "
		}

		paramIN += "$" + strconv.Itoa(idx+1)
		params = append(params, orgIds[idx])
	}

	result, err := e.tx.Query(e.ctx, `
								SELECT organization_id, category_id
								FROM organization_by_category
								WHERE organization_by_category.organization_id IN (`+paramIN+`)`, params...)
	if err != nil {
		return
	}
	defer result.Close()

	var (
		id       int64
		category int64
	)
	res = make(map[int64][]int64, len(orgIds))
	for result.Next() {
		if err = result.Scan(&id, &category); err != nil {
			break
		}
		res[id] = append(res[id], category)
	}

	return
}
