package parkingentity

import (
	"doit-parking/pkg/queuex"
	"fmt"
)

type (
	// Spot represents a parking spot in the parking lot.
	Spot struct {
		Floor int
		Col   int
		Row   int
	}

	// VehicleSpot represents a parked vehicle information.
	VehicleSpot struct {
		SpotID
		Type        VehicleType
		StillParked bool
	}
)

// AvailableSpots holds the available parking spots for different vehicle types.
type AvailableSpots struct {
	B1 *queuex.Queue[Spot]
	M1 *queuex.Queue[Spot]
	A1 *queuex.Queue[Spot]
}

// SpotID represents a unique identifier for a parking spot with format: floor-row-col.
type SpotID Spot

func (s SpotID) ID() string {
	return fmt.Sprintf("%d-%d-%d", s.Floor, s.Row, s.Col)
}

// VehicleType represents the type of vehicle that can be parked in the parking lot.
type VehicleType uint

const (
	// M1 represents a motorcycle.
	M1 VehicleType = iota
	// B1 represents a bicycles.
	B1
	// A1 represents an automobiles.
	A1
	// X0 represents an inactive parking spot.
	X0
)
