package parkingcli

import (
	"github.com/mtfiqh/DoiT-parking-system/parking/parkingentity"
	"github.com/mtfiqh/DoiT-parking-system/pkg/queuex"
)

func (p *parking) Unpark(spotID string, vehicleNumber int) error {
	p.mutex.RLock()
	vehicleSpot, exists := p.VehiclesParked[vehicleNumber]
	p.mutex.RUnlock()
	if !exists {
		return parkingentity.ErrVehicleNotFound
	}

	if vehicleSpot.SpotID.ID() != spotID {
		return parkingentity.ErrVehicleNotFound
	}

	if vehicleSpot.StillParked == false {
		return parkingentity.ErrVehicleNotFound
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	vehicleSpot.StillParked = false
	p.VehiclesParked[vehicleNumber] = vehicleSpot

	var qfunc *queuex.Queue[parkingentity.Spot]
	// Find an available spot for the vehicle type
	switch vehicleSpot.Type {
	case parkingentity.A1:
		qfunc = p.AvailableSpots.A1
	case parkingentity.B1:
		qfunc = p.AvailableSpots.B1
	case parkingentity.M1:
		qfunc = p.AvailableSpots.M1
	default:
		return parkingentity.ErrInvalidVehicleType
	}

	qfunc.Enqueue(parkingentity.Spot{
		Floor: vehicleSpot.Floor,
		Col:   vehicleSpot.Col,
		Row:   vehicleSpot.Row,
	})

	return nil
}
