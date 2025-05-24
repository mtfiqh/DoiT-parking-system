package cli

import (
	"doit-parking/cli/parkingcli"
	"log"
)

func RunParkingSimulation() {
	log.Println("Running parking simulation...")

	park, err := parkingcli.NewPark(parkingcli.WithRandomizeParkingSpots(8, 1000, 1000))
	if err != nil {
		log.Fatal(err)
	}

	_ = park
}
