package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/maruqu/alien-invasion/internal/simulation"
	"github.com/maruqu/alien-invasion/internal/world"
)

const (
	defaultIterationsLimit = 10000
	defaultAliensCount     = 50
)

var (
	iterationsLimit   int
	aliensCount       int
	outputMapFilepath string

	runCmd = &cobra.Command{
		Use:   "run [input map file]",
		Short: "Run simulation",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			log.SetFlags(0)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			worldMap, err := world.Load(args[0])
			if err != nil {
				return fmt.Errorf("error loading world map: %w", err)
			}

			simulation, err := simulation.NewSimulation(
				iterationsLimit,
				aliensCount,
				worldMap,
			)
			if err != nil {
				return fmt.Errorf("error initializing simulation: %w", err)
			}

			result, err := simulation.Run()
			if err != nil {
				return fmt.Errorf("error running simulation: %w", err)
			}

			if outputMapFilepath != "" {
				err = world.Save(outputMapFilepath, result)
				if err != nil {
					return fmt.Errorf("error saving result world map: %w", err)
				}
			} else {
				if len(result) == 0 {
					log.Println("Whole world destroyed!")
				} else {
					log.Printf("\nWorld map after invasion:\n\n%s", result)
				}
			}

			return nil
		},
	}
)

func init() {
	runCmd.Flags().IntVarP(&iterationsLimit, "iterations", "i", defaultIterationsLimit, "iterations limit")
	runCmd.Flags().IntVarP(&aliensCount, "aliens", "a", defaultAliensCount, "aliens count")
	runCmd.Flags().StringVarP(&outputMapFilepath, "output", "o", "", "output world map file (printed to STDOUT by default)")
}
