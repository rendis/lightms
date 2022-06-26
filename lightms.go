package lightms

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Run load properties, inject dependencies, run primaries and start lightms server
func Run() {
	// Populating dByAlias and dByTypes
	processInjectionsByReceivers()
	processInjectionsByDirect()

	// Resolve Properties
	resolveProps()

	// Resolve Dependencies
	resolveDependencies()

	// Run primaries
	runPrimaries()

	// Start lightms server
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	// waiting for the out signal
	s := <-signalChan
	log.Printf("Out signal triggered: %v", s)
	os.Exit(0)
}
