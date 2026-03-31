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
	MinAdvanceTime int     `gorm:"default:720" json:"min_advance_time"` // Minimum advance time in minutes (default 12h = 720m)
	MaxActiveBookings int  `gorm:"default:3" json:"max_active_bookings"` // Max active bookings per user per room
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
