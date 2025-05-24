package parking

import (
	"doit-parking/parking/parkingentity"
	"doit-parking/pkg/queuex"
	"sync"
)

// parking provides an implementation of a parking system that allows vehicles to be parked, unparked, and searched for within a structured parking space.
type parking struct {
	Spaces         [][][]int // floor, col, row
	AvailableSpots *parkingentity.AvailableSpots
	VehiclesParked map[int]parkingentity.VehicleSpot

	mutex *sync.RWMutex
}

// Parking interface defines the methods for parking operations.
type Parking interface {
	Park(vehicleType parkingentity.VehicleType, vehicleNumber int) (*parkingentity.SpotID, error)
	Unpark(spotID string, vehicleNumber int) error
	AvailableSpot(vehicleType parkingentity.VehicleType) (int, []parkingentity.Spot)
	SearchVehicle(vehicleNumber int) (*parkingentity.SpotID, error)
}

// NewPark initializes a new parking instance with empty spaces and available spots.
func NewPark(opts ...ParkOption) (Parking, error) {
	// get options
	opt := &ParkOptions{}
	for _, o := range opts {
		o(opt)
	}

	park := &parking{
		Spaces: make([][][]int, 0),
		AvailableSpots: &parkingentity.AvailableSpots{
			B1: queuex.NewQueue[parkingentity.Spot](),
			M1: queuex.NewQueue[parkingentity.Spot](),
			A1: queuex.NewQueue[parkingentity.Spot](),
		},
		VehiclesParked: make(map[int]parkingentity.VehicleSpot),
		mutex:          new(sync.RWMutex),
	}

	if opt.WithRandomize {
		err := park.Seed(opt.MaxFloor, opt.MaxCol, opt.MaxRow)
		if err != nil {
			return nil, err
		}
	}

	return park, nil
}

// ParkOptions defines the options for initializing a parking instance.
type ParkOptions struct {
	WithRandomize bool
	MaxFloor      int
	MaxCol        int
	MaxRow        int
}

// ParkOption is a function type that modifies the ParkOptions.
type ParkOption func(*ParkOptions)

// WithRandomizeParkingSpots is an option to initialize the parking with random spots (seeding).
func WithRandomizeParkingSpots(maxFloor, maxCol, maxRow int) ParkOption {
	return func(opt *ParkOptions) {
		opt.MaxFloor = maxFloor
		opt.MaxCol = maxCol
		opt.MaxRow = maxRow
		opt.WithRandomize = true
	}
}
