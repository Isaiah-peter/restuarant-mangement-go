package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `json:"_id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	OrderDate time.Time          `json:"order_date" validate:"required"`
	OrderId   string             `json:"order_id"`
	TableId   *string            `json:"table_id" validate:"required"`
}
