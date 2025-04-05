package models

type Dashboard struct {
	User          User
	UserEquipment []UserEquipment
	Tasks         []UserMaintenanceTask
}
