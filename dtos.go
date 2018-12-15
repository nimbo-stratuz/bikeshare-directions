package main

import (
	"time"
)

type userID string

// Bicycle ...
type Bicycle struct {
	ID            int       `json:"id"`
	SmartLockUUID string    `json:"smartLockUUID"`
	Available     bool      `json:"available"`
	DateAdded     time.Time `json:"dateAdded"`
	Location      Location  `json:"location"`
	OwnerID       userID    `json:"ownerId"`
	Rentals       []Rental  `json:"rentals"`
}

// Rental ...
type Rental struct {
	BicycleID     int       `json:"bicycleId"`
	RentStart     time.Time `json:"rentStart"`
	RentEnd       time.Time `json:"rentEnd"`
	StartLocation Location  `json:"startLocation"`
	EndLocation   Location  `json:"endLocation"`
	ID            int       `json:"id"`
	UserID        userID    `json:"userId"`
}

// Location ...
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
