package lightms

import (
	"log"
)

// PrimaryProcess is the primary process to be run by the lightms server
type PrimaryProcess interface {
	Start()
}

type primaryProcessSupply func() PrimaryProcess

type primary struct {
	name     string
	supplier primaryProcessSupply
}

var primaries = make([]*primary, 0)

// addPrimary adds a primary process to the list of primaries
func addPrimary(name string, pps primaryProcessSupply) {
	primaries = append(primaries, &primary{name, pps})
}

// runPrimaries runs all the primaries
func runPrimaries() {
	for _, pps := range primaries {
		pp := pps.supplier()
		log.Printf("Running primary process: %s\n", pps.name)
		go pp.Start()
	}
}
