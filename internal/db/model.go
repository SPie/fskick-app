package db

import (
	"time"

	"github.com/spie/fskick/internal/uuid"
)

type Model struct {
	ID        uint      `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"-"`
	UUID      string    `json:"uuid"`
}

func (model *Model) CreateUUID() error {
	uuid, err := uuid.GenerateUuidString()
	if err != nil {
		return err
	}

	model.UUID = uuid

	return nil
}
