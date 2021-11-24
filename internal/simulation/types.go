package simulation

import (
	"fmt"
	"strings"
)

// WorldMap is an internal representation of a parsed input file.
type WorldMap map[City]Neighbors

func (wm WorldMap) String() string {
	var sb strings.Builder

	for city, neighbors := range wm {
		parts := make([]string, 0, 5)
		parts = append(parts, string(city))

		if neighbors.North != "" {
			parts = append(parts, fmt.Sprintf("north=%s", neighbors.North))
		}
		if neighbors.South != "" {
			parts = append(parts, fmt.Sprintf("south=%s", neighbors.South))
		}
		if neighbors.East != "" {
			parts = append(parts, fmt.Sprintf("east=%s", neighbors.East))
		}
		if neighbors.West != "" {
			parts = append(parts, fmt.Sprintf("west=%s", neighbors.West))
		}

		sb.WriteString(strings.Join(parts, " ") + "\n")
	}

	return sb.String()
}

// Neighbors respresent connections to other cities.
// nil pointer represents no connection for a direction.
type Neighbors struct {
	North City
	South City
	East  City
	West  City
}

type City string

type AlienPositions map[Alien]City

type Alien string
