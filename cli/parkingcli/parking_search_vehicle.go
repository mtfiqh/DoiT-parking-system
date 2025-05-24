package parkingcli

import (
	parkingentity2 "doit-parking/parking/parkingentity"
)

func (p *parking) SearchVehicle(vehicleNumber int) (*parkingentity2.SpotID, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	vehicleSpot, exists := p.VehiclesParked[vehicleNumber]
	if !exists {
		return nil, parkingentity2.ErrVehicleNotFound
	}

	return &vehicleSpot.SpotID, nil
}
