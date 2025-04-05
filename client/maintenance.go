package client

import "github.com/noodahl-org/erp/internal/models"

func (c *erpClient) FetchUserMaintenanceTasks(query models.UserMaintenanceTask) ([]models.UserMaintenanceTask, error) {
	results := []models.UserMaintenanceTask{}
	res, err := c.client.R().
		SetResult(&results).
		Get("maintenance/user/" + query.UserID)
	if err != nil {
		return nil, err
	}
	return results, res.Err
}
