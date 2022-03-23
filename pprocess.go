package lightms

import (
	"log"
)

// PrimaryProcess is the primary process to be run by the lightms server
type PrimaryProcess interface{ Start() }

type primaryProcessFunc func() PrimaryProcess

var primaries = make([]func() PrimaryProcess, 0)

// AddPrimary adds a primary process to the list of primaries
func AddPrimary(primary primaryProcessFunc) {
	primaries = append(primaries, primary)
}

// runPrimaries runs all the primaries
func runPrimaries() {
	for _, primary := range primaries {
		pp := primary()
		log.Printf("Running primary process: %T\n", pp)
		go pp.Start()
	}
}
