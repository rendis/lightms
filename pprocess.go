package lightms

import (
	"log"
	"reflect"
)

var primaries []PrimaryProcess
var pptype = reflect.TypeOf((*PrimaryProcess)(nil)).Elem()

// PrimaryProcess is the primary process to be run by the lightms server
type PrimaryProcess interface {
	Start()
}

// runPrimaries runs all the primaries
func runPrimaries() {
	if len(primaries) == 0 {
		log.Fatalf("No primaries process found. Please add one or more primary processes to the lightms server.")
	}
	log.Printf("Running %d primaries process", len(primaries))
	for _, pps := range primaries {
		go pps.Start()
	}
}
