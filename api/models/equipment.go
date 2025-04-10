package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Equipment struct {
	ID          string         `json:"id" gorm:"type:uuid;primary_key;"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at"`
	Make        string         `json:"make" gorm:"uniqueIndex:equip_composite_idx;"`
	Model       string         `json:"model" gorm:"uniqueIndex:equip_composite_idx;"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	Tags        pq.StringArray `json:"tags" gorm:"type:text[]"`
}

func (e *Equipment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}

	return nil
}

func (e Equipment) ListLabel() string {
	return fmt.Sprintf("%s - %s", e.Make, e.Model)
}

type UserEquipment struct {
	ID           string     `json:"id" gorm:"type:uuid;primary_key;"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
	UserID       string     `json:"user_id" gorm:"type:uuid"`
	EquipmentID  string     `json:"equipment_id" gorm:"uuid;"`
	Equipment    Equipment  `json:"equipment" gorm:"foreignKey:EquipmentID"`
	SerialNumber string     `json:"serial_number"`
	Year         int        `json:"year"`
}

func (e *UserEquipment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}

	return nil
}

func (e UserEquipment) ListLabel() string {
	return fmt.Sprintf("%s %s", e.Equipment.Make, e.Equipment.Model)
}

type EquipmentComponent struct {
	ID          string     `json:"id" gorm:"type:uuid;primary_key;"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	EquipmentID string     `json:"equipment_id"`
	Name        string     `json:"name"`
}

func (e EquipmentComponent) ListLabel() string {
	return e.Name
}
