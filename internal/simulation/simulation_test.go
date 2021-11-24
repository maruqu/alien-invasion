package simulation

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	simpleMap = WorldMap{
		"Talihina": Neighbors{
			South: "Pinson",
		},
		"Pinson": Neighbors{
			North: "Talihina",
			East:  "Fabens",
		},
		"Fabens": Neighbors{
			West: "Pinson",
		},
		"Clifton": Neighbors{},
	}

	starMap = WorldMap{
		"Centercity": Neighbors{
			North: "Northcity",
			South: "Southcity",
			East:  "Eastcity",
			West:  "Westcity",
		},
		"Northcity": Neighbors{
			South: "Centercity",
		},
		"Southcity": Neighbors{
			North: "Centercity",
		},
		"Eastcity": Neighbors{
			West: "Centercity",
		},
		"Westcity": Neighbors{
			East: "Centercity",
		},
	}
)

func Test_NewSimulation(t *testing.T) {
	t.Run("aliens placed on map correctly", func(t *testing.T) {
		s, err := NewSimulation(10, 3, simpleMap)
		require.NoError(t, err)

		assert.Len(t, s.alienPositions, 3)

		for _, city := range s.alienPositions {
			assert.Contains(t, simpleMap, city)
		}
	})
}

func Test_Simulation(t *testing.T) {
	t.Run("aliens and city destroyed", func(t *testing.T) {
		s := &Simulation{
			iterationCounter: 0,
			worldMap:         copyMap(starMap),
			alienPositions: AlienPositions{
				"Alien 1": "Centercity",
				"Alien 2": "Centercity",
			},
		}

		s.Step()

		expectedMap := WorldMap{
			"Northcity": Neighbors{},
			"Southcity": Neighbors{},
			"Eastcity":  Neighbors{},
			"Westcity":  Neighbors{},
		}

		assert.Equal(t, 1, s.iterationCounter)
		assert.Equal(t, expectedMap, s.worldMap)
		assert.Empty(t, s.alienPositions)
	})

	t.Run("alien takes an existing road", func(t *testing.T) {
		s := &Simulation{
			iterationCounter: 0,
			worldMap:         copyMap(simpleMap),
			alienPositions: AlienPositions{
				"Alien 1": "Talihina",
			},
		}
		s.Step()

		assert.Equal(t, 1, s.iterationCounter)
		assert.Equal(t, City("Pinson"), s.alienPositions["Alien 1"])
	})

	t.Run("alien never visits an isolated city", func(t *testing.T) {
		s := &Simulation{
			iterationCounter: 0,
			iterationLimit:   100,
			worldMap:         copyMap(simpleMap),
			alienPositions: AlienPositions{
				"Alien 1": "Pinson",
			},
		}

		visitedCities := make(map[City]struct{})
		for i := 0; i < 100; i++ {
			s.Step()

			visitedCities[s.alienPositions["Alien 1"]] = struct{}{}
		}

		assert.NotContains(t, s.alienPositions["Alien 1"], City("Clifton"))
		assert.Equal(t, 100, s.iterationCounter)
	})

	t.Run("alien does not move when trapped in an isolated city", func(t *testing.T) {
		s := &Simulation{
			iterationCounter: 0,
			iterationLimit:   100,
			worldMap:         copyMap(simpleMap),
			alienPositions: AlienPositions{
				"Alien 1": "Clifton",
			},
		}

		s.Step()

		assert.Equal(t, City("Clifton"), s.alienPositions["Alien 1"])
	})

	t.Run("alien is able to travel in any valid direction", func(t *testing.T) {
		s := &Simulation{
			iterationCounter: 0,
			iterationLimit:   100,
			worldMap:         copyMap(starMap),
			alienPositions: AlienPositions{
				"Alien 1": "Centercity",
			},
		}

		visitedCities := make(map[City]struct{})
		for i := 0; i < 100; i++ {
			s.Step()

			visitedCities[s.alienPositions["Alien 1"]] = struct{}{}
		}

		assert.Len(t, visitedCities, 5)
	})

	t.Run("simulation runs until iteration limit is reached", func(t *testing.T) {
		s, err := NewSimulation(100, 1, copyMap(simpleMap))
		require.NoError(t, err)

		result, err := s.Run()
		require.NoError(t, err)

		assert.Equal(t, simpleMap, result)
		assert.Equal(t, 100, s.iterationCounter)
		assert.True(t, s.ShouldStop())
	})

	t.Run("names with ids generated if more that 75 aliens", func(t *testing.T) {
		s, err := NewSimulation(100, 76, copyMap(simpleMap))
		require.NoError(t, err)

		assert.Contains(t, s.alienPositions, Alien("Alien 1"))
	})
}
