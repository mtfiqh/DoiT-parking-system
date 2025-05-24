package parkingcli

import (
	"github.com/mtfiqh/DoiT-parking-system/parking/parkingentity"
)

func (p *parking) SearchVehicle(vehicleNumber int) (*parkingentity.SpotID, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	vehicleSpot, exists := p.VehiclesParked[vehicleNumber]
	if !exists {
		return nil, parkingentity.ErrVehicleNotFound
	}

	return &vehicleSpot.SpotID, nil
}
