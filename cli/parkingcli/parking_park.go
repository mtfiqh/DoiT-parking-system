package parkingcli

import (
	"doit-parking/parking/parkingentity"
	"doit-parking/pkg/queuex"
)

func (p *parking) Park(vehicleType parkingentity.VehicleType, vehicleNumber int) (*parkingentity.SpotID, error) {
	// Check if the vehicle is already parked
	p.mutex.RLock()
	parked, exists := p.VehiclesParked[vehicleNumber]
	p.mutex.RUnlock()
	if exists && parked.StillParked {
		return nil, parkingentity.ErrVehicleAlreadyParked
	}

	var qfunc *queuex.Queue[parkingentity.Spot]
	// Find an available spot for the vehicle type
	switch vehicleType {
	case parkingentity.A1:
		qfunc = p.AvailableSpots.A1
	case parkingentity.B1:
		qfunc = p.AvailableSpots.B1
	case parkingentity.M1:
		qfunc = p.AvailableSpots.M1
	default:
		return nil, parkingentity.ErrInvalidVehicleType
	}

	// Dequeue a spot from the available spots queue
	spot, ok := qfunc.Dequeue()
	if !ok {
		return nil, parkingentity.ErrSpotNotFound
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	spotID := parkingentity.SpotID{
		Floor: spot.Floor,
		Col:   spot.Col,
		Row:   spot.Row,
	}

	p.VehiclesParked[vehicleNumber] = parkingentity.VehicleSpot{
		SpotID:      spotID,
		Type:        vehicleType,
		StillParked: true,
	}

	return &spotID, nil
}
