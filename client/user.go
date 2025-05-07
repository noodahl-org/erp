package client

import "github.com/noodahl-org/erp/api/models"

func (c *erpClient) FetchUser(query models.User) (*models.User, error) {
	var user models.User
	res, err := c.client.R().
		SetResult(&user).
		SetQueryParams(map[string]string{
			"id":    query.ID,
			"make":  query.Username,
			"model": query.Password,
		}).
		Get("/users")
	if err != nil {
		return nil, err
	}
	return &user, res.Err
}

func (c *erpClient) FetchDashboard(userID string) (*models.Dashboard, error) {
	var dash models.Dashboard
	res, err := c.client.R().
		SetResult(&dash).
		SetQueryParams(map[string]string{
			"user_id": userID,
		}).
		Get("/dashboard")
	if err != nil {
		return nil, err
	}
	return &dash, res.Err
}
