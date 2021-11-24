package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/maruqu/alien-invasion/cmd"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
