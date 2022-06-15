package lightms

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

var instanceReg = newInstanceRegistry()
var confRr = newConfReader(instanceReg.addInstance, instanceReg.addDependent)

// Run load properties, inject dependencies, run primaries and start lightms server
func Run() {
	confRr.readConfReg()
	instanceReg.resolveDependents()
	runPrimaries()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	// waiting for the out signal
	s := <-signalChan
	log.Printf("Out signal triggered: %v", s)
	os.Exit(0)
}

func AddConf[T any](conf *T) {
	confRr.addConf(conf)
}
