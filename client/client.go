package client

import (
	"github.com/noodahl-org/erp/internal/models"
	"resty.dev/v3"
)

type ERPClient interface {
	FetchUser(query models.User) (*models.User, error)
	FetchEquipment(query models.Equipment) ([]models.Equipment, error)
	DeleteUserEquipment(id string) error
	UpsertEquipment(query models.Equipment) (*models.Equipment, error)
	FetchUserEquipment(query models.UserEquipment) ([]models.UserEquipment, error)
	UpserUserEquipment(query models.UserEquipment) (*models.UserEquipment, error)
	FetchUserMaintenanceTasks(query models.UserMaintenanceTask) ([]models.UserMaintenanceTask, error)
}

type erpClient struct {
	client *resty.Client
}

func NewERPClient(opts ...func(*erpClient)) ERPClient {
	cl := &erpClient{
		client: resty.New(),
	}
	cl.client.SetHeader("Content-Type", "application/json")
	for _, opt := range opts {
		opt(cl)
	}
	return cl
}

func WithBaseURL(url string) func(*erpClient) {
	return func(e *erpClient) {
		e.client.SetBaseURL(url)
	}
}
