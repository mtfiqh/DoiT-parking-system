package parking

import (
	"doit-parking/parking/parkingentity"
	"doit-parking/pkg/queuex"
)

func (p *parking) AvailableSpot(vehicleType parkingentity.VehicleType) (int, []parkingentity.Spot) {
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
		return 0, nil
	}

	return qfunc.Size, qfunc.Print()
}
