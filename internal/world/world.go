package world

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/maruqu/alien-invasion/internal/simulation"
)

// Load reads and parses a world map from a provided file.
// Map structure is not validated.
func Load(filepath string) (simulation.WorldMap, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// parse lines
	worldMap := make(simulation.WorldMap, len(lines))
	for _, line := range lines {
		neighbors := simulation.Neighbors{}

		parts := strings.Split(line, " ")
		if len(parts) < 1 || strings.Contains(parts[0], "=") {
			return nil, fmt.Errorf("error parsing map file: line without city")
		}

		for _, part := range parts[1:] {
			directionCity := strings.Split(part, "=")
			if len(directionCity) != 2 {
				return nil, fmt.Errorf("error parsing map file: invalid road: %s", part)
			}

			city := simulation.City(directionCity[1])

			switch directionCity[0] {
			case "north":
				neighbors.North = city
			case "south":
				neighbors.South = city
			case "east":
				neighbors.East = city
			case "west":
				neighbors.West = city
			}
		}

		worldMap[simulation.City(parts[0])] = neighbors
	}

	return worldMap, nil
}

// Save writes a world map to a provided filepath.
func Save(filepath string, worldMap simulation.WorldMap) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(worldMap.String())

	return nil
}
