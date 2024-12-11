package models

import "github.com/lib/pq"

type Room struct {
    RoomID         int            `gorm:"primaryKey" json:"room_id"`
    RatePerNight   float64        `json:"rate_per_night"`
    MaxGuests      int            `json:"max_guests"`
    AvailableDates pq.StringArray `gorm:"type:text[]" json:"available_dates"` // Use pq.StringArray for PostgreSQL arrays
}
