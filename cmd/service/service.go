package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/alexdemen/companies/app"
	"github.com/alexdemen/companies/app/models"

	"github.com/go-chi/chi/v5"
)

type CompanyService struct {
	CompanyLogic *app.Catalog
}

func (cs *CompanyService) AddBuilding(w http.ResponseWriter, r *http.Request) error {
	var building models.Building
	if err := readBody(r, &building); err != nil {
		return err
	}

	id, err := cs.CompanyLogic.AddBuilding(building)
	if err != nil {
		return err
	}

	return fillBody(w, IdResponse{id})
}

func (cs *CompanyService) GetOrganizationByIds(w http.ResponseWriter, r *http.Request) error {
	var ids []int64
	if err := readBody(r, &ids); err != nil {
		return err
	}

	res, err := cs.CompanyLogic.GetOrganizationByIds(ids)
	if err != nil {
		return err
	}

	return fillBody(w, res)
}

func (cs *CompanyService) GetOrganizationInBuilding(w http.ResponseWriter, r *http.Request) error {
	buildingId := chi.URLParam(r, "building_id")

	bId, err := strconv.ParseInt(buildingId, 10, 64)
	if err != nil {

	}

	res, err := cs.CompanyLogic.GetOrganizationInBuilding(bId)
	if err != nil {
		return err
	}

	return fillBody(w, res)
}

func (cs *CompanyService) GetOrganizationByCategory(w http.ResponseWriter, r *http.Request) error {
	categoryId := chi.URLParam(r, "category_id")

	cId, err := strconv.ParseInt(categoryId, 10, 64)
	if err != nil {
		return err
	}

	res, err := cs.CompanyLogic.GetOrganizationByCategory(cId)
	if err != nil {
		return err
	}

	return fillBody(w, res)
}

func fillBody(w http.ResponseWriter, data interface{}) error {
	resBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resBody)

	return err
}

func readBody(r *http.Request, dest interface{}) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.New("error read request body")
	}

	err = json.Unmarshal(body, dest)
	if err != nil {
		return errors.New("error unmarshal request body")
	}

	return nil
}
