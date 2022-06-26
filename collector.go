package lightms

import (
	"log"
	"reflect"
)

var props = make(map[reflect.Type]reflect.Value) // Use to store all loaded props

var completedPP []*container         // Used to set completed dependencies to call constructorFunc
var completedContainers []*container // Used to set completed dependencies to call constructorFunc

var isolateInstances []*container   // Instances with no dependencies
var dependentInstances []*container // Instances with dependencies

var dByAlias = make(map[string][]*container)       // dependencies by alias
var dByTypes = make(map[reflect.Type][]*container) // dependencies by type

var aliasResolvers = make(map[string]*container)        // used to resolve dependencies by alias
var typeResolvers = make(map[reflect.Type][]*container) // used to resolve dependencies by type

// addAliasResolver adds a resolver for the given alias.
func addAliasResolver(alias string, c *container) {
	if _, ok := aliasResolvers[alias]; ok {
		log.Fatalf("alias '%s' already exists", alias)
	}
	aliasResolvers[alias] = c
}

// addTypeResolver adds a resolver for the given type.
func addTypeResolver(t reflect.Type, c *container) {
	typeResolvers[t] = append(typeResolvers[t], c)
}

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
func notifyTypeResolved(val reflect.Value) {
	t := val.Type()
	if v, ok := dByTypes[t]; ok {
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
		delete(dByTypes, t)
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
		notifyTypeResolved(prop)
	}

	// While there are completed dependencies, resolve them
	for len(completedContainers) > 0 || len(completedPP) > 0 {
		runCompletedContainer()
		runCompletedPP()
	}

	// If there are still dependencies, then there is a circular dependency or a dependency is not resolved
	if len(dByAlias) > 0 || len(dByTypes) > 0 {
		log.Println("There is a circular dependency or a dependency is not resolved")
		if len(dByAlias) > 0 {
			log.Println("Dependencies by alias:")
			for k, v := range dByAlias {
				log.Printf("%s: %v", k, v)
			}
		}
		if len(dByTypes) > 0 {
			log.Println("Dependencies by type:")
			for k, v := range dByTypes {
				log.Printf("%s: %v", k, v)
			}
		}
		log.Fatalln("Dependencies not resolved")
	}
}

// runCompletedPP used to call constructorFunc of PrimaryProcess
func runCompletedPP() {
	for len(completedPP) > 0 {
		c := completedPP[0]
		completedPP = completedPP[1:]
		val := c.constructorFunc.Call(c.params)[0]
		primaries = append(primaries, val.Interface().(PrimaryProcess))
		alterVal := getAlter(val) //Get alter value, val is a pointer then alterVal is the value, otherwise alterVal is a pointer to val
		notifyTypeResolved(val)
		notifyTypeResolved(alterVal)
		for _, a := range c.aliases {
			notifyAliasResolved(a, val)
			notifyAliasResolved(a, alterVal)
		}
	}
}

// runCompletedContainer used to call constructorFunc of container
func runCompletedContainer() {
	for len(completedContainers) > 0 {
		c := completedContainers[0]
		completedContainers = completedContainers[1:]
		val := c.constructorFunc.Call(c.params)[0]
		alterVal := getAlter(val) //Get alter value, val is a pointer then alterVal is the value, otherwise alterVal is a pointer to val
		notifyTypeResolved(val)
		notifyTypeResolved(alterVal)
		for _, a := range c.aliases {
			notifyAliasResolved(a, val)
			notifyAliasResolved(a, alterVal)
		}
	}
}

// runIsolates used to call constructorFunc of isolateInstances
func runIsolates() {
	for _, c := range isolateInstances {
		val := c.constructorFunc.Call([]reflect.Value{})[0]
		alter := getAlter(val) //Get alter value, val is a pointer then alterVal is the value, otherwise alterVal is a pointer to val
		notifyTypeResolved(val)
		notifyTypeResolved(alter)
		for _, a := range c.aliases {
			notifyAliasResolved(a, val)
			notifyAliasResolved(a, alter)
		}
	}
}

// getAlter used to get alter value, val is a pointer then alterVal is the value, otherwise alterVal is a pointer to val
func getAlter(val reflect.Value) reflect.Value {
	if val.Type().Kind() == reflect.Ptr {
		return val.Elem()
	}
	return reflect.New(val.Type()).Elem()
}
