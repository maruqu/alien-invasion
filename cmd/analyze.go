package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"

	"github.com/maruqu/alien-invasion/internal/util"
	"github.com/maruqu/alien-invasion/internal/world"
)

var (
	analyzeCmd = &cobra.Command{
		Use:   "analyze [initial map file] [result map file] [initial dot file] [output dot file]",
		Short: "Generate a graph in dot format from the simulation result with destroyed cities marked red.",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			initialWorldMap, err := world.Load(args[0])
			if err != nil {
				return fmt.Errorf("error loading initial map: %w", err)
			}

			resultWorldMap, err := world.Load(args[1])
			if err != nil {
				return fmt.Errorf("error loading result map: %w", err)
			}

			var destroyedCities []string
			for city, _ := range initialWorldMap {
				if _, ok := resultWorldMap[city]; !ok {
					destroyedCities = append(destroyedCities, string(city))
				}
			}

			b, err := ioutil.ReadFile(args[2])
			if err != nil {
				return fmt.Errorf("error reading dot graph: %w", err)
			}
			graph := string(b)

			// add background color property to destroyed nodes
			for _, city := range destroyedCities {
				oldAttrs := fmt.Sprintf("[label=\"%s\"]", city)
				newAttrs := fmt.Sprintf("[label=\"%s\", fillcolor=\"red\"]", city)
				graph = strings.Replace(graph, oldAttrs, newAttrs, 1)
			}

			err = util.Write(args[3], graph)
			if err != nil {
				return fmt.Errorf("error writing generated dot graph to file: %w", err)
			}

			return nil
		},
	}
)
