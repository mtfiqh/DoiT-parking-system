package parkingentity

import "errors"

var (
	ErrVehicleAlreadyParked = errors.New("vehicle is already parked")
	ErrInvalidVehicleType   = errors.New("invalid vehicle type")
	ErrSpotNotFound         = errors.New("spot not found")
	ErrVehicleNotFound      = errors.New("vehicle not found")
)
