package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	UUID string `gorm:"unique;not null"`
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
