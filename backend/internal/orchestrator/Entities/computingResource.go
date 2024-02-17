package Entities

import (
	"216/internal/Types"
	"github.com/gofrs/uuid"
)

type ComputingResource struct {
	Id             uuid.UUID
	Name           string
	Task           *uuid.UUID
	TaskExpression *ArithmeticExpressions `gorm:"foreignkey:Task;" json:"task_expression"`
	TaskStr        *string
	Status         Types.Status
	HeartBeat      int
}
