package lightms

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// LoadLightMs load properties from yml file and run primary processes
func LoadLightMs() {
	loadProperties()
	runPrimaries()
}

// Run load properties, run primaries and start lightms server
func Run() {
	LoadLightMs()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	// waiting for the out signal
	s := <-signalChan
	log.Printf("Out signal triggered: %v", s)
	os.Exit(0)
}
