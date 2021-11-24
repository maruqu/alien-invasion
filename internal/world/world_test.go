package world

import (
	"path"
	"testing"

	"github.com/maruqu/alien-invasion/internal/simulation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testMap = simulation.WorldMap{
		"Talihina": simulation.Neighbors{
			South: "Pinson",
		},
		"Pinson": simulation.Neighbors{
			North: "Talihina",
			East:  "Fabens",
		},
		"Fabens": simulation.Neighbors{
			West: "Pinson",
		},
	}
)

func Test_Save_Load(t *testing.T) {
	tempDir := t.TempDir()
	filepath := path.Join(tempDir, "test.map")

	err := Save(filepath, testMap)
	require.NoError(t, err)

	loadedMap, err := Load(filepath)
	require.NoError(t, err)

	assert.EqualValues(t, testMap, loadedMap)
}
