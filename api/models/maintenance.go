package models

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type MaintenanceTask struct {
	ID          string     `json:"id" gorm:"type:uuid;primary_key;"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	EquipmentID string     `json:"equipment_id" gorm:"type:uuid;"`
	Equipment   Equipment  `gorm:"foreignKey:EquipmentID" json:"-"`
	Cron        string     `validate:"is_cron" json:"cron"`
	Description *string    `json:"description"`
	Priority    *string    `validate:"oneof=low medium high" json:"priority"`
}

func (m *MaintenanceTask) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	
	return nil
}

type UserMaintenanceTask struct {
	ID                string          `json:"id" gorm:"type:uuid;primary_key;"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         *time.Time      `json:"deleted_at"`
	StartTime         *time.Time      `json:"start_time"`
	UserID            string          `json:"user_id" gorm:"type:uuid"`
	Assignee          User            `gorm:"foreignKey:UserID" json:"-"`
	MaintenanceTaskID string          `json:"maintenance_task_id"`
	MaintenanceTask   MaintenanceTask `gorm:"foreignKey:MaintenanceTaskID" json:"maintenance_task"`
	Completed         *bool           `json:"completed"`
}

func (u *UserMaintenanceTask) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().String()
	return nil
}
