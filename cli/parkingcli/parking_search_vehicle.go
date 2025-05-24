package parkingcli

import (
	"doit-parking/parking/parkingentity"
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
