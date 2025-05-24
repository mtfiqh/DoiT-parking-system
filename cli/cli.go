package cli

import (
	"context"
	"doit-parking/cli/parkingcli"
	"doit-parking/parking/parkingentity"
	"doit-parking/pkg/randomizer"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"sync"
	"time"
)

func RunParkingSimulation(floor, column, row, gates int, duration time.Duration) error {

	log.Println("Running parking simulation...")
	log.Printf("floor: %d, column: %d, row: %d, gates: %d", floor, column, row, gates)

	park, err := parkingcli.NewPark(parkingcli.WithRandomizeParkingSpots(floor, column, row))
	if err != nil {
		log.Fatal(err)
	}

	// get initial available spots
	errg, _ := errgroup.WithContext(context.Background())
	var (
		a1countChan = make(chan int, 1)
		b1countChan = make(chan int, 1)
		m1countChan = make(chan int, 1)
	)
	errg.Go(func() error {
		a1Count, _ := park.AvailableSpot(parkingentity.A1)
		a1countChan <- a1Count
		return nil
	})

	errg.Go(func() error {
		b1Count, _ := park.AvailableSpot(parkingentity.B1)
		b1countChan <- b1Count
		return nil
	})

	errg.Go(func() error {
		m1Count, _ := park.AvailableSpot(parkingentity.M1)
		m1countChan <- m1Count
		return nil
	})

	if err := errg.Wait(); err != nil {
		log.Fatal(errors.Wrap(err, "error getting initial available spots"))
		return err
	}

	beforeA1 := <-a1countChan
	beforeB1 := <-b1countChan
	beforeM1 := <-m1countChan

	log.Printf("Initial spots - A1: %d, B1: %d, M1: %d", beforeA1, beforeB1, beforeM1)

	close(a1countChan)
	close(b1countChan)
	close(m1countChan)

	log.Println("Parking simulation start in 5s...")
	time.Sleep(5 * time.Second)

	tnow := time.Now().Local()
	tend := tnow.Add(duration)

	wg, _ := errgroup.WithContext(context.Background())

	wg.SetLimit(gates) // how many gates to simulate

	parked := make(map[int]parkingentity.SpotID)
	mu := new(sync.RWMutex)

	i := 0

	/**
	function to simulate parking operations
	with chance of
	- parking: 60%
	- unparking: 30%
	- searching: 10%
	 **/
	for time.Now().Local().Before(tend) {
		i++
		func(i int) {
			wg.Go(func() error {
				op := randomizer.RandomizeInt(1, 10)

				// parking operation
				if op <= 6 {
					vehicleNum := 10000 + i
					vehicleType := randomizer.RandomizeEnum(parkingentity.M1, parkingentity.B1, parkingentity.A1)

					spotID, err := park.Park(vehicleType, vehicleNum)
					if err != nil {
						switch errors.Cause(err) {
						case parkingentity.ErrSpotNotFound:
							// means full, do nothing
							log.Printf("Parking full for vehicle %d of type %d", vehicleNum, vehicleType)
						default:
							err = errors.Wrap(err, fmt.Sprintf("parking vehicle %d of type %d", vehicleNum, vehicleType))
							log.Println("Error parking vehicle:", err)
							return err
						}

						return nil
					}

					mu.Lock()
					defer mu.Unlock()
					parked[vehicleNum] = *spotID

					log.Printf("parked vehicle: %v, in: %v", vehicleNum, spotID.ID())
					return nil
				}

				// unparking operation
				if op <= 9 {
					mu.Lock()
					defer mu.Unlock()

					// if empty, unpark nothing
					if len(parked) == 0 {
						return nil
					}

					// get a random vehicle to unpark
					// todo optimize
					keys := make([]int, 0, len(parked))
					for k := range parked {
						keys = append(keys, k)
					}

					vehicleNum := randomizer.RandomizeEnum(keys...)

					// unpark method
					err := park.Unpark(parked[vehicleNum].ID(), vehicleNum)
					if err != nil {
						err = errors.Wrap(err, fmt.Sprintf("unparking vehicle %d", vehicleNum))
						log.Println("Error unparking vehicle:", err)
						return err
					}

					// delete in parked map
					delete(parked, vehicleNum)

					log.Println("Unparked vehicle:", vehicleNum)
					return nil
				}

				// searching operation
				if op <= 10 {
					mu.RLock()

					keys := make([]int, 0, len(parked))

					if len(parked) == 0 {
						mu.RUnlock()
						return nil
					}

					for k := range parked {
						keys = append(keys, k)
					}

					vehicleNum := randomizer.RandomizeEnum(keys...)
					mu.RUnlock()

					spotID, err := park.SearchVehicle(vehicleNum)
					if err != nil {
						err = errors.Wrap(err, fmt.Sprintf("searching vehicle %d", vehicleNum))
						log.Println("Error searching vehicle:", err)
						return err
					}

					if spotID == nil {
						return errors.New("spotID empty")
					} else {
						log.Printf("Vehicle %d found at spot %s", vehicleNum, spotID.ID())
					}

					return nil
				}

				return nil
			})

		}(i)
	}

	if err := wg.Wait(); err != nil {
		log.Fatal(errors.Wrap(err, "error in parking simulation"))
		return err
	}

	log.Println("Parking simulation completed successfully.")
	a1, _ := park.AvailableSpot(parkingentity.A1)
	b1, _ := park.AvailableSpot(parkingentity.B1)
	m1, _ := park.AvailableSpot(parkingentity.M1)
	log.Println("RESULT:")
	log.Printf("before: A1: %d, B1: %d, M1: %d", beforeA1, beforeB1, beforeM1)
	log.Printf("after: A1: %d, B1: %d, M1: %d", a1, b1, m1)
	log.Printf("remaining vehicles parked: %d", len(parked))

	totalBefore := beforeA1 + beforeB1 + beforeM1
	totalAfter := a1 + b1 + m1
	log.Printf("total spots: %d, total free spots: %d, remaining + free spots: %d", totalBefore, totalAfter, totalAfter+len(parked))
	log.Printf("total executions: %d", i)

	return nil
}
