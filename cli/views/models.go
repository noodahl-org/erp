package views

import (
	"github.com/noodahl-org/erp/api/models"
	"github.com/noodahl-org/erp/client"
)

type ViewOpts struct {
	ERP  client.ERPClient
	User *models.User
}
