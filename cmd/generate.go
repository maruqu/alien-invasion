package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/maruqu/alien-invasion/internal/mapgen"
	"github.com/maruqu/alien-invasion/internal/util"
)

const (
	defaultGridHeight  = 5
	defaultGridWidth   = 5
	defaultCitiesCount = 20
)

var (
	gridHeight       int
	gridWidth        int
	citiesCount      int
	dotGraphFilepath string

	generateCmd = &cobra.Command{
		Use:   "generate [output map file]",
		Short: "Generate a world map",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			gridMap, err := mapgen.NewGridMap(gridHeight, gridWidth, citiesCount)
			if err != nil {
				return fmt.Errorf("error generating map: %w", err)
			}

			err = util.Write(args[0], gridMap.String())
			if err != nil {
				return fmt.Errorf("error writing generated map to file: %w", err)
			}

			if dotGraphFilepath != "" {
				graph, err := gridMap.DotGraph()
				if err != nil {
					return fmt.Errorf("error generating dot format graph: %w", err)
				}

				err = util.Write(dotGraphFilepath, graph)
				if err != nil {
					return fmt.Errorf("error writing generated dot graph to file: %w", err)
				}
			}

			return nil
		},
	}
)

func init() {
	generateCmd.Flags().IntVarP(&gridHeight, "height", "", defaultGridHeight, "grid height")
	generateCmd.Flags().IntVarP(&gridWidth, "width", "", defaultGridWidth, "grid width")
	generateCmd.Flags().IntVarP(&citiesCount, "cities", "c", defaultCitiesCount, "cities count")
	generateCmd.Flags().StringVarP(&dotGraphFilepath, "dot", "d", "", "output dot file (graphviz format)")
}
