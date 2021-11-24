package mapgen

import (
	"bytes"
	_ "embed"
	"fmt"
	"math/rand"
	"strings"
	"text/template"
)

//go:embed city-names.txt
var cityNames string

//go:embed grid.dot.tmpl
var dotGraphTemplate string

// GridMap stores a generated map.
// Two internal representations are used to store information about city coordinates.
type GridMap struct {
	grid     [][]string
	worldMap map[string]neighbors
}

type neighbors struct {
	north *city
	south *city
	east  *city
	west  *city
}

type city struct {
	name        string
	coordinates coordinates
}

type coordinates [2]int

// NewGrid returns initialized GridMap structure containing a world map generated using provided parameters.
// The process of generation is following:
// 1. A grid of size height x width is created.
// 2. Provided number of cities is randomly placed on the grid.
// 3. If two cities are in the same row or column and there are no other cities between them, a road is created.
func NewGridMap(height, width, citiesCount int) (*GridMap, error) {
	if height*width < citiesCount {
		return nil, fmt.Errorf(
			"error creating grid: too many cities (%d) for provided map dimensions (%dx%d)",
			citiesCount, height, width,
		)
	}

	grid, err := generateGrid(height, width, citiesCount)
	if err != nil {
		return nil, err
	}

	worldMap := generateWorldMap(grid)

	return &GridMap{
		grid:     grid,
		worldMap: worldMap,
	}, nil
}

func (gm *GridMap) String() string {
	var sb strings.Builder

	for city, neighbors := range gm.worldMap {
		sb.WriteString(city)

		if neighbors.north != nil {
			sb.WriteString(fmt.Sprintf(" north=%s", neighbors.north.name))
		}

		if neighbors.south != nil {
			sb.WriteString(fmt.Sprintf(" south=%s", neighbors.south.name))
		}

		if neighbors.east != nil {
			sb.WriteString(fmt.Sprintf(" east=%s", neighbors.east.name))
		}

		if neighbors.west != nil {
			sb.WriteString(fmt.Sprintf(" west=%s", neighbors.west.name))
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

// DotGraph generates a dot format graph representation of world map.
// Dot language is used by Graphviz (https://graphviz.org).
// TODO: refactor
func (gm *GridMap) DotGraph() (string, error) {
	var sb strings.Builder

	// generate vertical grid structure
	column := make([]string, len(gm.grid))
	for i := 0; i < len(gm.grid[0]); i++ {
		for j := 0; j < len(gm.grid); j++ {
			column[j] = fmt.Sprintf("N%d_%d", j, i)
		}

		sb.WriteString(strings.Join(column, " -- ") + "\n")
	}

	verticalEdges := sb.String()
	sb.Reset()

	// generate horizontal grid structure
	row := make([]string, len(gm.grid[0]))
	for i := 0; i < len(gm.grid); i++ {
		for j := 0; j < len(gm.grid[0]); j++ {
			row[j] = fmt.Sprintf("N%d_%d", i, j)
		}

		sb.WriteString(fmt.Sprintf("rank=same {%s}\n", strings.Join(row, " -- ")))
	}

	horizontalEdges := sb.String()
	sb.Reset()

	// hide nodes without cities
	for i := 0; i < len(gm.grid); i++ {
		for j := 0; j < len(gm.grid[0]); j++ {
			if gm.grid[i][j] == "" {
				sb.WriteString(fmt.Sprintf("N%d_%d [style=invis]\n", i, j))
			}
		}
	}

	hiddenNodes := sb.String()
	sb.Reset()

	// draw roads between cities
	connectedCities := make(map[string]struct{})
	for i := 0; i < len(gm.grid); i++ {
		for j := 0; j < len(gm.grid[0]); j++ {
			if gm.grid[i][j] != "" {
				neighbors := gm.worldMap[gm.grid[i][j]]

				if neighbors.north != nil {
					if _, ok := connectedCities[neighbors.north.name]; !ok {
						sb.WriteString(fmt.Sprintf("N%d_%d -- N%d_%d [style=solid]\n", i, j, neighbors.north.coordinates[0], neighbors.north.coordinates[1]))
					}
				}
				if neighbors.south != nil {
					if _, ok := connectedCities[neighbors.south.name]; !ok {
						sb.WriteString(fmt.Sprintf("N%d_%d -- N%d_%d [style=solid]\n", i, j, neighbors.south.coordinates[0], neighbors.south.coordinates[1]))
					}
				}
				if neighbors.east != nil {
					if _, ok := connectedCities[neighbors.east.name]; !ok {
						sb.WriteString(fmt.Sprintf("N%d_%d -- N%d_%d [style=solid]\n", i, j, neighbors.east.coordinates[0], neighbors.east.coordinates[1]))
					}
				}
				if neighbors.west != nil {
					if _, ok := connectedCities[neighbors.west.name]; !ok {
						sb.WriteString(fmt.Sprintf("N%d_%d -- N%d_%d [style=solid]\n", i, j, neighbors.west.coordinates[0], neighbors.west.coordinates[1]))
					}
				}

				connectedCities[gm.grid[i][j]] = struct{}{}
			}
		}
	}

	roads := sb.String()
	sb.Reset()

	// label nodes with city names
	for i := 0; i < len(gm.grid); i++ {
		for j := 0; j < len(gm.grid[0]); j++ {
			if gm.grid[i][j] != "" {
				sb.WriteString(fmt.Sprintf("N%d_%d [label=\"%s\"]\n", i, j, gm.grid[i][j]))
			}
		}
	}

	labels := sb.String()

	// insert generated nodes and edges into the template

	tmpl, err := template.New("").Parse(dotGraphTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, map[string]string{
		"hiddenNodes":     hiddenNodes,
		"verticalEdges":   verticalEdges,
		"horizontalEdges": horizontalEdges,
		"roads":           roads,
		"labels":          labels,
	})

	if err != nil {
		return "", fmt.Errorf("error generating graph from template: %s", err)
	}

	// TODO destroyed cities placeholder

	return result.String(), nil
}

// generateGrid returns a grid of size height x width with cities placed in random places.
func generateGrid(height, width, citiesCount int) ([][]string, error) {
	grid := make([][]string, height)
	for i := range grid {
		grid[i] = make([]string, width)
	}

	// place cities on the grid

	cities, err := getCityNames(citiesCount)
	if err != nil {
		return nil, err
	}

	for _, city := range cities {
		// Pick a random location until an empty location is found.
		// This is not the most efficient method and can take many iterations when
		// there is a small number of empty spots left.
		for {
			h, w := rand.Intn(height), rand.Intn(width)

			// ensure that another city is not already placed here
			if grid[h][w] != "" {
				continue
			}

			grid[h][w] = city
			break
		}
	}

	return grid, nil
}

// generateWorldMap creates a mapping of city names to its neighbors based on the provided grid.
func generateWorldMap(grid [][]string) map[string]neighbors {
	worldMap := make(map[string]neighbors)

	// get all possible roads from each city
	for h := 0; h < len(grid); h++ {
		for w := 0; w < len(grid[0]); w++ {
			if grid[h][w] != "" {
				neighbors := findNeighbors(h, w, grid)
				worldMap[grid[h][w]] = neighbors
			}
		}
	}

	return worldMap
}

// getCityNames returns a slice of cities with a provided count.
func getCityNames(count int) ([]string, error) {
	names := strings.Split(strings.TrimSpace(cityNames), "\n")

	if len(names) < count {
		return nil, fmt.Errorf("maximum number of cities exceeded (%d)", len(names))
	}

	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = names[i]
	}

	return result, nil
}

// findNeighbors finds closest cities in the same row or column of the grid.
func findNeighbors(h, w int, grid [][]string) neighbors {
	result := neighbors{}

	// north
	for i := h - 1; i >= 0; i-- {
		if grid[i][w] != "" {
			result.north = &city{
				name:        grid[i][w],
				coordinates: coordinates{i, w},
			}
			break
		}
	}

	// south
	for i := h + 1; i < len(grid); i++ {
		if grid[i][w] != "" {
			result.south = &city{
				name:        grid[i][w],
				coordinates: coordinates{i, w},
			}
			break
		}
	}

	// east
	for i := w + 1; i < len(grid[0]); i++ {
		if grid[h][i] != "" {
			result.east = &city{
				name:        grid[h][i],
				coordinates: coordinates{h, i},
			}
			break
		}
	}

	// west
	for i := w - 1; i >= 0; i-- {
		if grid[h][i] != "" {
			result.west = &city{
				name:        grid[h][i],
				coordinates: coordinates{h, i},
			}
			break
		}
	}

	return result
}
