package parking

import (
	"doit-parking/parking/parkingentity"
)

// Parking interface defines the methods for parking operations.
type Parking interface {
	Park(vehicleType parkingentity.VehicleType, vehicleNumber int) (*parkingentity.SpotID, error)
	Unpark(spotID string, vehicleNumber int) error
	AvailableSpot(vehicleType parkingentity.VehicleType) (int, []parkingentity.Spot)
	SearchVehicle(vehicleNumber int) (*parkingentity.SpotID, error)
}
