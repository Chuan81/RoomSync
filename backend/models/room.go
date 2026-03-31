package models

import (
	"time"
)

type Room struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Capacity    int       `json:"capacity"`
	Location    string    `json:"location"`
	Equipment   string    `json:"equipment"` // e.g., "Projector, Whiteboard"
	Status      string    `gorm:"default:'available'" json:"status"` // 'available', 'maintenance'
	NeedApproval bool      `gorm:"default:false" json:"need_approval"` // true if admin approval is required
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
