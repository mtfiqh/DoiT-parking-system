package parkingcli

import (
	"context"
	"github.com/mtfiqh/DoiT-parking-system/parking/parkingentity"
	"github.com/mtfiqh/DoiT-parking-system/pkg/randomizer"
	"golang.org/x/sync/errgroup"
)

func (p *parking) Seed(maxFloor, maxCol, maxRow int) error {
	wg, _ := errgroup.WithContext(context.Background())
	wg.SetLimit(10)

	p.Spaces = make([][][]int, maxFloor)
	for i := 0; i < maxFloor; i++ {
		p.Spaces[i] = make([][]int, maxRow)
		for j := 0; j < maxRow; j++ {
			p.Spaces[i][j] = make([]int, maxCol)
		}
	}

	for i := 0; i < maxFloor; i++ {
		for j := 0; j < maxRow; j++ {
			for k := 0; k < maxCol; k++ {
				floor, row, col := i, j, k
				spot := randomizer.RandomizeEnum(parkingentity.A1, parkingentity.B1, parkingentity.M1, parkingentity.X0)

				p.Spaces[floor][row][col] = int(spot)
				switch spot {
				case parkingentity.B1:
					p.AvailableSpots.B1.Enqueue(parkingentity.Spot{Floor: floor, Col: col, Row: row})
				case parkingentity.M1:
					p.AvailableSpots.M1.Enqueue(parkingentity.Spot{Floor: floor, Col: col, Row: row})
				case parkingentity.A1:
					p.AvailableSpots.A1.Enqueue(parkingentity.Spot{Floor: floor, Col: col, Row: row})
				default:
					// do nothing
				}

			}
		}
	}

	return nil
}
