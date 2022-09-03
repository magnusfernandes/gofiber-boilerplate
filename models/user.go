package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id            uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name          string          `gorm:"column:name" json:"name"`
	Phone         string          `gorm:"column:phone;unique" json:"phone"`
	Email         string          `gorm:"column:email;unique" json:"email"`
	Role          string          `gorm:"column:role;default:end_user" json:"role"`
	ProfilePic    string          `gorm:"column:profile_pic" json:"profilePic"`
	PasswordHash  string          `gorm:"column:password_hash" json:"passwordHash"`
	Gender        string          `gorm:"column:gender" json:"gender"`
	BirthDate     sql.NullTime    `gorm:"column:birth_date" json:"birthDate"`
}

func (u User) ShortJson() map[string]interface{} {
	payload := map[string]interface{}{
		"id":    u.Id,
		"name":  u.Name,
		"email": u.Email,
		"role":  u.Role,
	}
	return payload
}

func (u User) Json() map[string]interface{} {
	payload := map[string]interface{}{
		"id":         u.Id,
		"name":       u.Name,
		"email":      u.Email,
		"phone":      u.Phone,
		"role":       u.Role,
		"profilePic": nil,
		"gender":     nil,
		"birthDate":  nil,
		"createdAt":  u.CreatedAt.Format(time.RFC3339),
	}
	if len(u.ProfilePic) > 0 {
		payload["profilePic"] = u.ProfilePic
	}
	if len(u.Gender) > 0 {
		payload["gender"] = u.Gender
	}
	if u.BirthDate.Valid {
		payload["birthDate"] = u.BirthDate.Time.Format(time.RFC3339)
	}
	return payload
}
