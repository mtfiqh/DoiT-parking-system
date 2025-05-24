package parkingcli

import (
	"context"
	"doit-parking/parking/parkingentity"
	"doit-parking/pkg/randomizer"
	"golang.org/x/sync/errgroup"
)

func (p *parking) Seed(maxFloor, maxCol, maxRow int) error {
	wg, _ := errgroup.WithContext(context.Background())
	wg.SetLimit(10)

	p.Spaces = make([][][]int, maxFloor)
	for i := 0; i < maxFloor; i++ {
		p.Spaces[i] = make([][]int, maxCol)
		for j := 0; j < maxCol; j++ {
			p.Spaces[i][j] = make([]int, maxRow)
		}
	}

	for i := 0; i < maxFloor; i++ {
		for j := 0; j < maxCol; j++ {
			for k := 0; k < maxRow; k++ {
				floor, col, row := i, j, k
				spot := randomizer.RandomizeEnum(parkingentity.A1, parkingentity.B1, parkingentity.M1, parkingentity.X0)

				p.Spaces[floor][col][row] = int(spot)
				switch spot {
				case parkingentity.B1:
					p.AvailableSpots.B1.Enqueue(parkingentity.Spot{Floor: floor, Col: col, Row: row})
				case parkingentity.M1:
					p.AvailableSpots.M1.Enqueue(parkingentity.Spot{Floor: floor, Col: col, Row: row})
				case parkingentity.A1:
					p.AvailableSpots.A1.Enqueue(parkingentity.Spot{Floor: floor, Col: col, Row: row})
				default:
					
				}

			}
		}
	}

	return nil
}
