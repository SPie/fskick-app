package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UUID      string         `gorm:"unique;not null" json:"uuid"`
}

func (model *Model) BeforeCreate(db *gorm.DB) error {
	p, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	model.UUID = p.String()

	return nil
}

type Repository interface {
	AutoMigrate()
}
