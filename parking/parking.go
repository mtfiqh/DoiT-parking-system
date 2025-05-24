package parking

import (
	"github.com/mtfiqh/DoiT-parking-system/parking/parkingentity"
)

// ParkingSystem interface defines the methods for parking operations.
type ParkingSystem interface {
	Park(vehicleType parkingentity.VehicleType, vehicleNumber int) (*parkingentity.SpotID, error)
	Unpark(spotID string, vehicleNumber int) error
	AvailableSpot(vehicleType parkingentity.VehicleType) (int, []parkingentity.Spot)
	SearchVehicle(vehicleNumber int) (*parkingentity.SpotID, error)
}
