package cmd

import (
	"fmt"
	"github.com/mtfiqh/DoiT-parking-system/cli"
	"github.com/spf13/cobra"
	"time"
)

var (
	gates    int
	floor    int
	rows     int
	column   int
	duration time.Duration
)

var simulateCmd = &cobra.Command{
	Use:   "cli:simulate",
	Short: "Simulate parking lot behavior",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚗 Simulating parking system with:")
		fmt.Printf("Gates  : %d\n", gates)
		fmt.Printf("Floors : %d\n", floor)
		fmt.Printf("Rows   : %d\n", rows)
		fmt.Printf("Columns: %d\n", column)
		fmt.Printf("Duration: %v\n", duration.String())

		// You can run your simulation logic here
		err := cli.RunParkingSimulation(floor, column, rows, gates, duration)
		if err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(simulateCmd)

	simulateCmd.Flags().IntVar(&gates, "gates", 10, "Number of gates to simulate multiple gates (park and unpark at the same time)")
	simulateCmd.Flags().IntVar(&floor, "floor", 8, "Number of floors")
	simulateCmd.Flags().IntVar(&rows, "rows", 1000, "Number of column per floor")
	simulateCmd.Flags().IntVar(&column, "column", 1000, "Number of columns per row")
	simulateCmd.Flags().DurationVar(&duration, "duration", 15*time.Second, "Duration of simulation")
}
