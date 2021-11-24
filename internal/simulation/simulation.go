package simulation

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	_ "embed"
)

//go:embed alien-names.txt
var alienNames string

// Simulation stores the state of the simulation.
type Simulation struct {
	iterationCounter int
	iterationLimit   int

	// citiesMap represents a graph as an adjacency list.
	worldMap WorldMap

	// alienPositions maps alien name to its current position (city).
	alienPositions map[Alien]City
}

// NewSimulation returned initialized Simulation structure with aliens randomly placed on the map.
func NewSimulation(iterationLimit, aliensCount int, worldMap WorldMap) (*Simulation, error) {
	if len(worldMap) == 0 {
		return nil, fmt.Errorf("map cannot be empty")
	}

	alienPositions := generateAlienPlacement(aliensCount, worldMap)

	s := &Simulation{
		iterationCounter: 0,
		iterationLimit:   iterationLimit,
		worldMap:         copyMap(worldMap),
		alienPositions:   alienPositions,
	}

	return s, nil
}

// Run starts simulation, executes steps until the stop condition is met and returns a WorldMap as result.
func (s *Simulation) Run() (WorldMap, error) {
	log.Println("Alien invasion started!")

	for !s.ShouldStop() {
		s.Step()
	}

	log.Println("Alien invasion finished!")

	return s.worldMap, nil
}

// ShouldStop returns true if a stop condition is met.
func (s *Simulation) ShouldStop() bool {
	return s.iterationCounter >= s.iterationLimit || len(s.alienPositions) == 0 || len(s.worldMap) == 0
}

// Step moves all aliens on the map and evaluate the rules.
func (s *Simulation) Step() {
	// evaluate the rules for the initial alien placement
	if s.iterationCounter == 0 {
		s.evaluateRules()
	}

	s.updateAlienPositions()
	s.evaluateRules()
	s.iterationCounter += 1
}

// generateAlienPlacement randomly assigns positions on the map for the provided alien count.
func generateAlienPlacement(aliensCount int, worldMap WorldMap) AlienPositions {
	aliens := getAliens(aliensCount)

	cities := make([]City, 0, len(worldMap))
	for city := range worldMap {
		cities = append(cities, city)
	}

	alienPositions := make(AlienPositions, len(aliens))
	for _, alien := range aliens {
		randomCityIdx := rand.Intn(len(cities))
		alienPositions[alien] = cities[randomCityIdx]
	}

	return alienPositions
}

// getAliens returns a slice of aliens with a provided count.
// Pre-generated names aer returned for up to 75 aliens.
// For greater counts, aliens are named ["Alien 1", "Alien 2",...]
func getAliens(count int) []Alien {
	result := make([]Alien, count)

	if count <= 75 {
		names := strings.Split(strings.TrimSpace(alienNames), "\n")

		for i := 0; i < count; i++ {
			result[i] = Alien(names[i])
		}
	} else {
		for i := 0; i < count; i++ {
			result[i] = Alien(fmt.Sprintf("Alien %d", i+1))
		}
	}

	return result
}

// updateAlienPositions calculates updated alien positions using connections between the cities.
func (s *Simulation) updateAlienPositions() {
	updatedAlienPositions := make(AlienPositions)

	for alien, city := range s.alienPositions {
		cityNeighbors := s.worldMap[city]

		// pick random direction
		possibleDirections := make([]City, 0, 4)
		if cityNeighbors.North != "" {
			possibleDirections = append(possibleDirections, cityNeighbors.North)
		}
		if cityNeighbors.South != "" {
			possibleDirections = append(possibleDirections, cityNeighbors.South)
		}
		if cityNeighbors.East != "" {
			possibleDirections = append(possibleDirections, cityNeighbors.East)
		}
		if cityNeighbors.West != "" {
			possibleDirections = append(possibleDirections, cityNeighbors.West)
		}

		// check if alien is trapped
		if len(possibleDirections) == 0 {
			updatedAlienPositions[alien] = city
			continue
		}

		randomIdx := rand.Intn(len(possibleDirections))
		updatedAlienPositions[alien] = possibleDirections[randomIdx]
	}

	s.alienPositions = updatedAlienPositions
}

// evaluateRules check for cities where two or more aliens are currently located in.
// Such cities and aliens are deleted from the simulation state.
func (s *Simulation) evaluateRules() {
	// check for alien fights
	cityAlienCount := make(map[City]int)
	for _, city := range s.alienPositions {
		cityAlienCount[city] += 1
	}

	// destroy aliens and cities
	for city, count := range cityAlienCount {
		if count >= 2 {
			// delete aliens
			var destroyedAlienNames []string
			for alien, position := range s.alienPositions {
				if position == city {
					destroyedAlienNames = append(destroyedAlienNames, string(alien))
					delete(s.alienPositions, alien)
				}
			}

			// delete city
			delete(s.worldMap, city)

			// delete all roads leading to this city
			for c, neighbors := range s.worldMap {
				if neighbors.North != "" && neighbors.North == city {
					neighbors.North = ""
				}
				if neighbors.South != "" && neighbors.South == city {
					neighbors.South = ""
				}
				if neighbors.East != "" && neighbors.East == city {
					neighbors.East = ""
				}
				if neighbors.West != "" && neighbors.West == city {
					neighbors.West = ""
				}

				s.worldMap[c] = neighbors
			}

			log.Printf(
				"%s has been destroyed by %s and %s!",
				city,
				strings.Join(destroyedAlienNames[:len(destroyedAlienNames)-1], ", "),
				destroyedAlienNames[len(destroyedAlienNames)-1],
			)
		}
	}
}

func copyMap(worldMap WorldMap) WorldMap {
	newWorldMap := make(WorldMap, len(worldMap))
	for c, n := range worldMap {
		newWorldMap[c] = n
	}
	return newWorldMap
}
