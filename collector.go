package lightms

import (
	"log"
	"reflect"
)

var props = make(map[reflect.Type]reflect.Value) // Use to store all loaded props

var isolateInstances []*container    // Instances with no dependencies waiting to call constructorFunc
var completedPP []*container         // Completed Primary Process waiting to call constructorFunc
var completedContainers []*container // Completed Containers waiting to call constructorFunc

var dependentInstances []*container // Instances with dependencies

var dByAlias = make(map[string][]*container)       // dependencies by alias
var dByTypes = make(map[reflect.Type][]*container) // dependencies by type

// addDependency adds not resolved dependency
func addDependency(c *container) {
	// If dependency not need parameter injection, then added to isolateInstances
	if c.paramInjections.Len() == 0 {
		isolateInstances = append(isolateInstances, c)
		return
	}
	dependentInstances = append(dependentInstances, c)
}

// processInjectionsByDirect analyze dependentInstances to populate dByAlias and dByTypes
func processInjectionsByDirect() {
	for _, c := range dependentInstances {
		for index, key := range c.paramInjections.Keys() {
			// If index == key (both are int), then it is a dependency by type
			if index == key {
				t, _ := c.paramInjections.Get(key)
				dByTypes[t] = append(dByTypes[t], c)
				continue
			}
			// If index != key (index(int), key(string)), then it is a dependency by alias
			s := key.(string)
			dByAlias[s] = append(dByAlias[s], c)
		}
	}
}

// storePropResolved saves the value of the property in the props map
func storePropResolved(prop reflect.Value) {
	if _, ok := props[prop.Type()]; !ok {
		props[prop.Type()] = prop
	}
}

// notifyTypeResolved used when a type is resolved to notify all dependencies of the same type
func notifyTypeResolved(typ reflect.Type, val reflect.Value) {
	if v, ok := dByTypes[typ]; ok {
		for _, c := range v {
			// completeInjectionByType returns true if all dependencies are resolved
			if c.completeInjectionByType(val) {
				// If dependency is PrimaryProcess, then added to completedPP
				if c.isPP {
					completedPP = append(completedPP, c)
				} else {
					completedContainers = append(completedContainers, c)
				}
			}
		}
		// Remove dependencies from dByTypes
		delete(dByTypes, typ)
	}
}

// notifyAliasResolved used when an alias is resolved to notify all dependencies of the same alias
func notifyAliasResolved(alias string, val reflect.Value) {
	if v, ok := dByAlias[alias]; ok {
		for _, c := range v {
			// completeInjectionByAlias returns true if all dependencies are resolved
			if c.completeInjectionByAlias(alias, val) {
				// If dependency is PrimaryProcess, then added to completedPP
				if c.isPP {
					completedPP = append(completedPP, c)
				} else {
					completedContainers = append(completedContainers, c)
				}
			}
		}
		// Remove dependencies from dByTypes
		delete(dByAlias, alias)
	}
}

// resolveDependencies resolve dependencies by alias and type
func resolveDependencies() {
	// First resolve isolateInstances
	runIsolates()

	// Notify props
	for _, prop := range props {
		notifyTypeResolved(prop.Type(), prop)
	}

	// While there are completed dependencies, resolve them
	for len(completedContainers) > 0 || len(completedPP) > 0 {
		runCompletedContainer()
		runCompletedPP()
	}

	// If there are still dependencies, then there is a circular dependency or a dependency is not resolved
	if len(dByAlias) > 0 || len(dByTypes) > 0 {
		log.Println("There is a circular dependency or a dependency is not resolved.")

		if len(dByAlias) > 0 {
			log.Println("Dependencies not resolved by alias:")
			for k, v := range dByAlias {
				log.Printf("- %s:", k)
				for _, c := range v {
					log.Printf("    - %v\n", c.constructorFunc.String())
				}
			}
		}

		if len(dByTypes) > 0 {
			log.Println("Dependencies not resolved by type:")
			for k, v := range dByTypes {
				log.Printf("- %s:", k)
				for _, c := range v {
					log.Printf("    - %v\n", c.constructorFunc.String())
				}
			}
		}
		log.Fatalln("Please check your injection configuration.")
	}
}

// runCompletedPP used to call constructorFunc of PrimaryProcess
func runCompletedPP() {
	pp := runCompleted(&completedPP)
	for _, c := range pp {
		addPP(c)
	}
}

// runCompletedContainer used to call constructorFunc of container
func runCompletedContainer() {
	runCompleted(&completedContainers)
}

// runIsolates used to call constructorFunc of isolateInstances
func runIsolates() {
	runCompleted(&isolateInstances)
}

func runCompleted(readyContainers *[]*container) []any {
	var callResults []any
	for len(*readyContainers) > 0 {
		r := *readyContainers
		c := r[0]
		*readyContainers = r[1:]
		val := c.constructorFunc.Call(c.params)[0]
		callResults = append(callResults, val.Interface())
		notifyResolved(c, val)
	}
	return callResults
}

// notifyResolved notify for each alias and interface in container
func notifyResolved(c *container, value reflect.Value) {
	for alias := range c.aliases {
		notifyAliasResolved(alias, value)
	}

	for interfaceTyp := range c.interfaces {
		notifyTypeResolved(interfaceTyp, value)
	}
}

// getAlter used to get alter value, val is a pointer then alterVal is the value, otherwise alterVal is a pointer to val
func getAlter(val reflect.Type) reflect.Type {
	if val.Kind() == reflect.Ptr {
		return val.Elem()
	}
	return reflect.PtrTo(val)
}
