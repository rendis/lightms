package lightms

import (
	"log"
)

// PrimaryProcess is the primary process to be run by the lightms server
type PrimaryProcess interface {
	Start()
}

type primaryProcessSupply func() PrimaryProcess

var primaries = make([]primaryProcessSupply, 0)

// addPrimary adds a primary process to the list of primaries
func addPrimary(primary primaryProcessSupply) {
	primaries = append(primaries, primary)
}

// runPrimaries runs all the primaries
func runPrimaries() {
	for _, primary := range primaries {
		pp := primary()
		log.Printf("Running primary process: %I\n", pp)
		go pp.Start()
	}
}
