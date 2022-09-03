package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OTPRequest struct {
	gorm.Model
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Code      string    `gorm:"column:code" json:"code"`
	IsExpired bool      `gorm:"column:is_expired;default:false" json:"isExpired"`
	UserId    uuid.UUID `gorm:"type:uuid;column:user_id" json:"userId"`
	User      User      `gorm:"foreignKey:UserId"`
}

func (OTPRequest) TableName() string {
	return "otp_request"
}

func (o OTPRequest) Json() map[string]interface{} {
	payload := map[string]interface{}{
		"id":        o.Id,
		"code":      o.Code,
		"createdAt": o.CreatedAt.Format(time.RFC3339),
	}
	return payload
}
