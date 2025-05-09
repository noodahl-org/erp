package client

import "github.com/noodahl-org/erp/api/models"

func (c *erpClient) FetchUserMaintenanceTasks(query models.UserMaintenanceTask) ([]models.UserMaintenanceTask, error) {
	results := []models.UserMaintenanceTask{}
	res, err := c.client.R().
		SetResult(&results).
		SetQueryParams(map[string]string{
			"id":           query.ID,
			"user_id":      query.UserID,
			"equipment_id": query.MaintenanceTask.EquipmentID,
		}).
		Get("maintenance/user/" + query.UserID)
	if err != nil {
		return nil, err
	}
	return results, res.Err
}

func (c *erpClient) FetchMaintenanceTasks(query models.MaintenanceTask) ([]models.MaintenanceTask, error) {
	results := []models.MaintenanceTask{}
	res, err := c.client.R().
		SetResult(&results).
		SetQueryParams(map[string]string{
			"id":           query.ID,
			"equipment_id": query.EquipmentID,
		}).
		Get("maintenance")
	if err != nil {
		return nil, err
	}
	return results, res.Err
}

func (c *erpClient) UpsertMaintenanceTask(query models.MaintenanceTask) (*models.MaintenanceTask, error) {
	result := models.MaintenanceTask{}
	res, err := c.client.R().
		SetBody(query).
		SetResult(&result).
		Post("maintenance")
	if err != nil {
		return nil, err
	}
	return &result, res.Err
}
