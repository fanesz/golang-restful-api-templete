package user

import (
	"backend/app/common/consts"
	"backend/app/common/model"

	"github.com/google/uuid"
)

type User struct {
	UUID         uuid.UUID `gorm:"primaryKey;type:uuid;not null" validate:"required,uuid4"`
	Username     string    `gorm:"uniqueIndex;size:100;not null"`
	Email        string    `gorm:"uniqueIndex;size:100;not null"`
	Password     string    `gorm:"size:255;not null"`
	RoleID       string    `gorm:"size:24;not null"`
	Role         Role      `gorm:"foreignKey:RoleID;references:ID"`
	LoginToken   string    `gorm:"size:255"`
	IPAddress    string    `gorm:"size:128"`
	ResetPWToken string    `gorm:"size:255"`
	model.TimestampsSoftDelete
}

type Role struct {
	ID   string      `gorm:"primaryKey;size:24;not null"`
	Name consts.Role `gorm:"uniqueIndex;size:100;not null" json:"role_name" validate:"required"`
	model.Timestamp
}

type UserGetByUUID struct {
	UUID string `uri:"uuid" validate:"required"`
}

type UserCreate struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserLogin struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
}

type UserUpdate struct {
	UUID        string      `uri:"uuid"`
	Username    string      `json:"username" validate:"min=3,max=100"`
	Email       string      `json:"email" validate:"email"`
	Password    string      `json:"password"`
	OldPassword string      `json:"old_password"`
	Role        consts.Role `json:"role" validate:"required"`
}

type UserResponse struct {
	UUID     string `json:"uuid,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role,omitempty"`
}

type CheckEmail struct {
	Email string `json:"email" validate:"required,email"`
}

type CheckUsername struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
}

type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordReset struct {
	NewPassword string `json:"new_password" validate:"required"`
}
