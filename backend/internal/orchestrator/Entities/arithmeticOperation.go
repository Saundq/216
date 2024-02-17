package Entities

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type ArithmeticOperation struct {
	ID       uuid.UUID `gorm:"primary_key;type:uuid;" json:"id"`
	Value    string    `gorm:"size:2; not null;type:string" json:"value"`
	LeadTime int       `json:"lead_time"`
}

func (ao *ArithmeticOperation) ChangeLeadTime(duration int) {
	ao.LeadTime = duration
}

func (ao *ArithmeticOperation) BeforeCreate(tx *gorm.DB) (err error) {
	uid, _ := uuid.NewV4()
	ao.ID = uid

	return
}
