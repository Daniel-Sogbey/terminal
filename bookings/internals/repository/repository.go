package repository

import (
	"bookings/internals/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res *models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
}
