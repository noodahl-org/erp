package client

import (
	"github.com/noodahl-org/erp/api/models"
)

func (c *erpClient) FetchEquipment(query models.Equipment) ([]models.Equipment, error) {
	results := []models.Equipment{}
	res, err := c.client.R().
		SetResult(&results).
		SetQueryParams(map[string]string{
			"id":    query.ID,
			"make":  query.Make,
			"model": query.Model,
		}).
		Get("/equipment")
	if err != nil {
		return nil, err
	}
	return results, res.Err
}

func (c *erpClient) UpsertEquipment(query models.Equipment) (*models.Equipment, error) {
	result := models.Equipment{}
	res, err := c.client.R().
		SetBody(query).
		SetResult(&result).
		Post("/equipment")
	if err != nil {
		return nil, err
	}
	return &result, res.Err
}

func (c *erpClient) FetchUserEquipment(query models.UserEquipment) ([]models.UserEquipment, error) {
	results := []models.UserEquipment{}
	res, err := c.client.R().
		SetResult(&results).
		Get("/equipment/user/" + query.UserID)
	if err != nil {
		return nil, err
	}
	return results, res.Err
}

func (c *erpClient) UpserUserEquipment(query models.UserEquipment) (*models.UserEquipment, error) {
	result := models.UserEquipment{}
	res, err := c.client.R().
		SetBody(query).
		SetResult(&result).
		Post("/equipment/user")
	if err != nil {
		return nil, err
	}
	return &result, res.Err
}

func (c *erpClient) DeleteUserEquipment(id string) error {
	res, err := c.client.R().
		Delete("/equipment/user/" + id)
	if err != nil {
		return err
	}
	return res.Err
}
