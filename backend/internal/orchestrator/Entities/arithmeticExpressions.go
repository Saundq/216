package Entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Status string

const (
	SUCCESS  Status = "Success"
	ERROR    Status = "Error"
	PROGRESS Status = "Progress"
	WHAIT    Status = "Wait"
)

type ArithmeticExpressions struct {
	ID                 uuid.UUID               `gorm:"primary_key;type:uuid;"         json:"id"`
	ExpressionString   string                  `gorm:"size:255; not null;type:string" json:"expression_string"`
	Status             Status                  `gorm:"size:10; not null;type:string"  json:"status"`
	AddedAt            int                     `gorm:"column:added_at;autoCreateTime" json:"added_at"`
	FinishedAt         int                     `json:"finished_at"`
	UpdatedAt          int                     `json:"updated_at"`
	Operand1           float64                 `json:"operand_1"`
	Operand2           float64                 `json:"operand_2"`
	Result             string                  `gorm:"size:255;type:string"           json:"result"`
	OrderField         int                     `json:"order_field"`
	StrValue           string                  `json:"str_value"`
	Operation          string                  `json:"operation"`
	Parent             *uuid.UUID              `gorm:"type:uuid;"                     json:"parent"`
	RequestId          *string                 `gorm:"uniqueIndex"                    json:"request_id"`
	Previous           *uuid.UUID              `gorm:"type:uuid;"                     json:"previous"`
	PreviousExpression *ArithmeticExpressions  `gorm:"foreignkey:Previous;"                 json:"previous_expression"`
	Next               *uuid.UUID              `gorm:"type:uuid;"                     json:"next"`
	NextExpression     *ArithmeticExpressions  `gorm:"foreignkey:Next;"                 json:"next_expression"`
	ExpressionPart     []ArithmeticExpressions `gorm:"foreignkey:Parent"`
	Owner              *uuid.UUID              `gorm:"type:uuid;"`
	//OwnerExpression    User                    `gorm:"foreignkey:ID"`
}

func (ae *ArithmeticExpressions) BeforeCreate(tx *gorm.DB) (err error) {
	//uid, _ := uuid.NewV4()
	//ae.ID = uid
	ae.Status = WHAIT

	return
}

func (ae *ArithmeticExpressions) SetResult(result string) {
	ae.Result = result
	ae.FinishedAt = int(time.Now().Unix())
	ae.Status = SUCCESS
}
