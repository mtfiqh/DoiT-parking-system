package parkingcli

import (
	parkingpkg "doit-parking/parking"
	"doit-parking/parking/parkingentity"
	"testing"
)

type parkingForDebug interface {
	parkingpkg.Parking
	GetSpaces() [][][]int
	GetAvailableSpots() *parkingentity.AvailableSpots
	GetVehiclesParked() map[int]parkingentity.VehicleSpot
}

func newParkForTest(opts ...ParkOption) (parkingForDebug, error) {
	p, err := NewPark(opts...)
	if err != nil {
		return nil, err
	}

	return p.(*parking), nil
}

func (p *parking) GetSpaces() [][][]int {
	return p.Spaces
}

func (p *parking) GetAvailableSpots() *parkingentity.AvailableSpots {
	return p.AvailableSpots
}

func (p *parking) GetVehiclesParked() map[int]parkingentity.VehicleSpot {
	return p.VehiclesParked
}

func TestSeedParkingSpots(t *testing.T) {
	const (
		maxFloors = 8
		maxRows   = 1000
		maxCols   = 1000
	)
	// Arrange
	p, err := newParkForTest(WithRandomizeParkingSpots(maxFloors, maxCols, maxRows))
	if err != nil {
		t.Fatalf("Failed to create parking: %v", err)
	}

	total := maxFloors * maxRows * maxCols

	countA1 := 0
	countB1 := 0
	countM1 := 0
	countX0 := 0

	spaces := p.GetSpaces()

	for i := 0; i < maxFloors; i++ {
		for j := 0; j < maxCols; j++ {
			for k := 0; k < maxRows; k++ {
				switch spaces[i][j][k] {
				case int(parkingentity.M1):
					countM1++
				case int(parkingentity.B1):
					countB1++
				case int(parkingentity.A1):
					countA1++
				case int(parkingentity.X0):
					countX0++
				}
			}
		}
	}

	// Assert
	if countA1+countB1+countM1+countX0 != total {
		t.Errorf("Total parking spots mismatch: expected %d, got %d", total, countA1+countB1+countM1+countX0)
	}

	avSpots := p.GetAvailableSpots()

	if countA1 != avSpots.A1.Size {
		t.Fatalf("Available A1 spots mismatch: expected %d, got %d", countA1, avSpots.A1.Size)
	}
	if countB1 != avSpots.B1.Size {
		t.Fatalf("Available B1 spots mismatch: expected %d, got %d", countB1, avSpots.B1.Size)
	}
	if countM1 != avSpots.M1.Size {
		t.Fatalf("Available M1 spots mismatch: expected %d, got %d", countM1, avSpots.M1.Size)
	}

	//	assertion each spots
	for i := 0; i < maxFloors; i++ {
		for j := 0; j < maxCols; j++ {
			for k := 0; k < maxRows; k++ {
				switch spaces[i][j][k] {
				case int(parkingentity.M1):
					if f, ok := avSpots.M1.Dequeue(); ok {
						if f.Floor != i || f.Row != k || f.Col != j {
							t.Fatalf("Expected M1 spot at (%d, %d, %d), got (%d, %d, %d)", i, j, k, f.Floor, f.Row, f.Col)
						}
					} else {
						t.Fatalf("Expected M1 spot to be available, but it was not")
					}

				case int(parkingentity.B1):
					if f, ok := avSpots.B1.Dequeue(); ok {
						if f.Floor != i || f.Row != k || f.Col != j {
							t.Fatalf("Expected B1 spot at (%d, %d, %d), got (%d, %d, %d)", i, j, k, f.Floor, f.Row, f.Col)
						}
					} else {
						t.Fatalf("Expected B1 spot to be available, but it was not")
					}
				case int(parkingentity.A1):
					if f, ok := avSpots.A1.Dequeue(); ok {
						if f.Floor != i || f.Row != k || f.Col != j {
							t.Fatalf("Expected A1 spot at (%d, %d, %d), got (%d, %d, %d)", i, j, k, f.Floor, f.Col, f.Row)
						}
					} else {
						t.Fatalf("Expected A1 spot to be available, but it was not")
					}
				}
			}
		}
	}
}

func TestPark(t *testing.T) {
	const (
		maxFloors = 5
		maxRows   = 100
		maxCols   = 100
	)

	t.Run("test parking", func(t *testing.T) {
		park, err := newParkForTest(WithRandomizeParkingSpots(maxFloors, maxCols, maxRows))
		if err != nil {
			t.Fatal(err)
		}

		spaces := park.GetSpaces()

		testCases := []struct {
			name          string
			vehicleNumber int
			vehicleType   parkingentity.VehicleType
			expectedError bool
		}{
			{
				name:          "parking A1",
				vehicleNumber: 1001,
				vehicleType:   parkingentity.A1,
				expectedError: false,
			},
			{
				name:          "parking A1 already park",
				vehicleNumber: 1001,
				vehicleType:   parkingentity.A1,
				expectedError: true,
			},
			{
				name:          "parking B1",
				vehicleNumber: 1010,
				vehicleType:   parkingentity.B1,
				expectedError: false,
			},
			{
				name:          "parking M1",
				vehicleNumber: 1100,
				vehicleType:   parkingentity.M1,
				expectedError: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				spotID, err := park.Park(tc.vehicleType, tc.vehicleNumber)
				if (err != nil) != tc.expectedError {
					t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
				}

				if err == nil && spotID == nil {
					t.Error("Expected a valid SpotID, got nil")
				}

				//	validate spotID
				if err == nil && spaces[spotID.Floor][spotID.Col][spotID.Row] != int(tc.vehicleType) {
					t.Error("Expected spot to be occupied by the vehicle type, but it was not")
				}

			})
		}

	})

	t.Run("test parking until space full", func(t *testing.T) {
		park, err := newParkForTest(WithRandomizeParkingSpots(maxFloors, maxCols, maxRows))
		if err != nil {
			t.Fatal(err)
		}

		ftest := func(vehicleType parkingentity.VehicleType) {
			currentTotalVehicleParked := len(park.GetVehiclesParked())
			totalAvailableSpaces := park.GetAvailableSpots().A1.Size

			for i := 0; i < totalAvailableSpaces; i++ {
				spotID, err := park.Park(vehicleType, 1000+i)

				if err != nil {
					t.Errorf("Failed to park %v vehicle %d: %v", vehicleType, 1000+i, err)
				}

				if spotID == nil {
					t.Error("Expected a valid SpotID, got nil")
				}
			}

			// Check if all A1 spots are occupied and try 1 more
			_, err := park.Park(vehicleType, 1000+totalAvailableSpaces)
			if err == nil {
				t.Errorf("Expected an error when parking %v, but got none", vehicleType)
			}

			if len(park.GetVehiclesParked())-currentTotalVehicleParked != totalAvailableSpaces {
				t.Errorf("Expected %d vehicles parked, got %d on parking %v", totalAvailableSpaces, len(park.GetVehiclesParked())-currentTotalVehicleParked, vehicleType)
			}
		}

		t.Run("parking A1", func(t *testing.T) {
			ftest(parkingentity.A1)
		})

		t.Run("parking B1", func(t *testing.T) {
			ftest(parkingentity.B1)
		})

		t.Run("parking M1", func(t *testing.T) {
			ftest(parkingentity.M1)
		})

	})

	t.Run("test park and unpark", func(t *testing.T) {
		park, err := newParkForTest(WithRandomizeParkingSpots(maxFloors, maxCols, maxRows))
		if err != nil {
			t.Fatal(err)
		}

		t.Run("unpark A1 vehicle not found", func(t *testing.T) {
			err := park.Unpark("1-1-1", 1001)
			if err == nil {
				t.Error("Expected an error when unparking a vehicle that is not parked, but got none")
			}
		})

		t.Run("unpark A1 vehicle found but spotID not match", func(t *testing.T) {
			_, err := park.Park(parkingentity.A1, 2001)
			if err != nil {
				t.Fatalf("Failed to park A1 vehicle: %v", err)
			}

			err = park.Unpark("0-0-0", 2001)
			if err == nil {
				t.Error("Expected an error when unparking with a mismatched spotID, but got none")
			}
		})

		t.Run("unpark A1, validate state park back to available spaces", func(t *testing.T) {
			spotID, err := park.Park(parkingentity.A1, 2002)
			if err != nil {
				t.Fatalf("Failed to park A1 vehicle: %v", err)
			}

			err = park.Unpark(spotID.ID(), 2002)
			if err != nil {
				t.Error("Expected no error when unparking a vehicle, but got:", err)
			}

			// Check if the spot is now available
			avSpots := park.GetAvailableSpots().A1.Print()
			avSpot := avSpots[len(avSpots)-1]

			spotLastID := parkingentity.SpotID{
				Floor: avSpot.Floor,
				Col:   avSpot.Col,
				Row:   avSpot.Row,
			}

			if spotLastID.ID() != spotID.ID() {
				t.Errorf("Expected the last available spot to be %s, got %s", spotID.ID(), spotLastID.ID())
			}

			t.Logf("Unparked vehicle %d from spot %s successfully, now spots %v available", 1002, spotID.ID(), spotLastID.ID())
		})
	})

	t.Run("test search vehicle", func(t *testing.T) {
		park, err := newParkForTest(WithRandomizeParkingSpots(maxFloors, maxCols, maxRows))
		if err != nil {
			t.Fatal(err)
		}

		t.Run("search vehicle not found", func(t *testing.T) {
			_, err := park.SearchVehicle(9999)
			if err == nil {
				t.Error("Expected an error when searching for a vehicle that is not parked, but got none")
			}
		})

		t.Run("search parked vehicle", func(t *testing.T) {
			parkSpotID, err := park.Park(parkingentity.A1, 1001)
			if err != nil {
				t.Fatalf("Failed to park A1 vehicle: %v", err)
			}

			spotID, err := park.SearchVehicle(1001)
			if err != nil {
				t.Fatalf("Failed to search parked vehicle: %v", err)
			}

			if spotID == nil {
				t.Error("Expected a valid SpotID, got nil")
			}

			if parkSpotID.ID() != spotID.ID() {
				t.Errorf("Expected SpotID %s, got %s", parkSpotID.ID(), spotID.ID())
			}
		})

		t.Run("search parked after unpark should be return last park", func(t *testing.T) {
			parkSpotID, err := park.Park(parkingentity.A1, 1003)
			if err != nil {
				t.Fatalf("Failed to park A1 vehicle: %v", err)
			}

			err = park.Unpark(parkSpotID.ID(), 1003)
			if err != nil {
				t.Fatalf("Failed to unpark A1 vehicle: %v", err)
			}

			spotID, err := park.SearchVehicle(1003)
			if err != nil {
				t.Error("Expected no error when searching for a parked vehicle, but got:", err)
			}

			if parkSpotID.ID() != spotID.ID() {
				t.Errorf("Expected SpotID %s, got %s", parkSpotID.ID(), spotID.ID())
			}
		})

		t.Run("search parked vehicle after park, unpark, park, unpark should be return last park", func(t *testing.T) {
			parkSpotID, err := park.Park(parkingentity.A1, 1004)
			if err != nil {
				t.Fatalf("Failed to park A1 vehicle: %v", err)
			}

			err = park.Unpark(parkSpotID.ID(), 1004)
			if err != nil {
				t.Fatalf("Failed to unpark A1 vehicle: %v", err)
			}

			parkSpotID, err = park.Park(parkingentity.A1, 1004)
			if err != nil {
				t.Fatalf("Failed to park A1 vehicle: %v", err)
			}

			err = park.Unpark(parkSpotID.ID(), 1004)
			if err != nil {
				t.Fatalf("Failed to unpark A1 vehicle: %v", err)
			}

			spotID, err := park.SearchVehicle(1004)
			if err != nil {
				t.Error("Expected no error when searching for a parked vehicle, but got:", err)
			}

			if parkSpotID.ID() != spotID.ID() {
				t.Errorf("Expected SpotID %s, got %s", parkSpotID.ID(), spotID.ID())
			}
		})
	})

	park, err := newParkForTest(WithRandomizeParkingSpots(maxFloors, maxCols, maxRows))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("test available spaces", func(t *testing.T) {
		checkFunc := func(vehicleType parkingentity.VehicleType) {
			total, spaces := park.AvailableSpot(vehicleType)
			allSpaces := park.GetSpaces()
			countSpace := 0

			for i, space := range allSpaces {
				for j, col := range space {
					for k, row := range col {
						if row == int(vehicleType) {
							if spaces[countSpace].Floor != i || spaces[countSpace].Col != j || spaces[countSpace].Row != k {
								t.Fatalf("Expected %v spot at (%d, %d, %d), got (%d, %d, %d)", vehicleType, i, j, k, spaces[countSpace].Floor, spaces[countSpace].Col, spaces[countSpace].Row)
							}

							countSpace++
						}
					}
				}
			}

			if total != countSpace {
				t.Fatalf("Expected %d available %v spaces, got %d", vehicleType, total, countSpace)
			}
		}

		t.Run("check A1 available spaces", func(t *testing.T) {
			checkFunc(parkingentity.A1)
		})

		t.Run("check B1 available spaces", func(t *testing.T) {
			checkFunc(parkingentity.B1)
		})

		t.Run("check M1 available spaces", func(t *testing.T) {
			checkFunc(parkingentity.M1)
		})

		t.Run("park 1 vehicle A1, spaces should be decreased", func(t *testing.T) {
			currentTotal, currentSpaces := park.AvailableSpot(parkingentity.A1)
			if currentTotal != len(currentSpaces) {
				t.Fatalf("Expected %d available A1 spaces before parking, got %d", currentTotal, len(currentSpaces))
			}

			spotID, err := park.Park(parkingentity.A1, 1005)
			if err != nil {
				t.Fatalf("Failed to park A1 vehicle: %v", err)
			}

			total, spaces := park.AvailableSpot(parkingentity.A1)
			if total != len(spaces) {
				t.Errorf("Expected %d available A1 spaces after parking, got %d", total, len(spaces))
			}

			if currentTotal-1 != total {
				t.Errorf("Expected %d available A1 spaces after parking, got %d", currentTotal-1, total)
			}

			err = park.Unpark(spotID.ID(), 1005)
			if err != nil {
				t.Fatalf("Failed to unpark A1 vehicle: %v", err)
			}

			total, spaces = park.AvailableSpot(parkingentity.A1)
			if total != len(spaces) {
				t.Errorf("Expected %d available A1 spaces after unparking, got %d", currentTotal, len(spaces))
			}

			if currentTotal != total {
				t.Errorf("Expected %d available A1 spaces after unparking, got %d", currentTotal, total)
			}

		})

	})
}
