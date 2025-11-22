package core

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type BaseModel struct {
	ID        uint64                `gorm:"primarykey"`
	DeletedAt soft_delete.DeletedAt `gorm:"index"`
	CreatedAt time.Time             `gorm:"Column:created_at"`
	CreatedBy string                `gorm:"Column:created_by"`
	UpdatedAt time.Time             `gorm:"Column:updated_at"`
	UpdatedBy string                `gorm:"Column:updated_by"`
}

// autofil audit fields
func (b *BaseModel) BeforeCreate(tx *gorm.DB) {
	b.CreatedBy = "1"
	b.UpdatedBy = "1"

}

func (b *BaseModel) BeforeUpdate(tx *gorm.DB) {
	b.UpdatedBy = "1"
}
