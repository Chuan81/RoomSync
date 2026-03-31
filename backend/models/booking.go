package models

import (
	"time"
)

type Booking struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	RoomID    uint      `gorm:"index;not null" json:"room_id"`
	Room      Room      `gorm:"foreignKey:RoomID" json:"room"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	StartTime time.Time `gorm:"index;not null" json:"start_time"`
	EndTime   time.Time `gorm:"index;not null" json:"end_time"`
	Title     string    `gorm:"not null" json:"title"`
	Status    string    `gorm:"default:'approved'" json:"status"` // 'pending', 'approved', 'rejected', 'cancelled', 'checked_in', 'completed', 'expired'
	CheckedIn bool      `gorm:"default:false" json:"checked_in"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
