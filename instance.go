package lightms

import (
	"github.com/rendis/orderedmap"
	"log"
	"reflect"
	"strings"
)

var pptype = reflect.TypeOf((*PrimaryProcess)(nil)).Elem()

type InstanceAlias interface {
	WithAlias(alias string, typ *reflect.Type) InstanceAlias
	AndInjections() InstanceInjection
}

type InstanceInjection interface {
	WithInjection(name string, typ *reflect.Type) InstanceInjection
}

func TOF[T any]() *reflect.Type {
	r := reflect.TypeOf((*T)(nil)).Elem()
	return &r
}

// AddInstance adds an instance to the container
func AddInstance(f any) InstanceAlias {
	return registerInstance(f, nil)
}

// AddInstanceWithAlias adds an instance to the container with an alias
func AddInstanceWithAlias(alias string, f any) InstanceAlias {
	return registerInstance(f, &alias)
}

func registerInstance(f any, alias *string) InstanceAlias {
	v := reflect.ValueOf(f)
	t := v.Type()
	checkConstructor(t)
	out := t.Out(0)

	// check if is a Primary Process
	isPP := false
	if out.Implements(pptype) {
		isPP = true
	}

	c := &container{
		constructorFunc: v,
		outType:         out,
		isPP:            isPP,
		aliases:         make(map[string]struct{}),
		interfaces:      make(map[reflect.Type]struct{}),
		paramsCount:     t.NumIn(),
		params:          make([]reflect.Value, t.NumIn()),
		paramInjections: orderedmap.New[reflect.Type](),
	}

	if alias != nil {
		al := transformAlias(*alias)
		c.aliases[al] = struct{}{}
	}

	// Registering the params to inject
	for i := 0; i < t.NumIn(); i++ {
		c.paramInjections.Set(i, t.In(i))
	}

	// Registering the instance
	addDependency(c)

	// First interface is the same as the outType
	c.interfaces[out] = struct{}{}
	if out.Kind() != reflect.Interface {
		alter := getAlter(out)
		// Adding the alter for the new type to the instance
		c.interfaces[alter] = struct{}{}
	}

	return c
}

type container struct {
	constructorFunc reflect.Value                       // constructor function
	outType         reflect.Type                        // output type
	isPP            bool                                // is Primary Process
	aliases         map[string]struct{}                 // aliases of the instance
	interfaces      map[reflect.Type]struct{}           // interfaces implemented by the instance
	params          []reflect.Value                     // input parameters to pass to the constructor
	paramsCount     int                                 // number of remaining parameters to pass to the constructor
	paramInjections orderedmap.OrderedMap[reflect.Type] // ordered map used to know how param are injected by alias or type
}

// WithAlias adds an alias with a type to the instance
func (c *container) WithAlias(alias string, interfaceTypePtr *reflect.Type) InstanceAlias {
	interfaceType := *interfaceTypePtr

	// interfaceType must be an interface
	if interfaceType.Kind() != reflect.Interface {
		log.Fatalf("[ERROR] Alias must be of type interface. Alias: %s, Type: %s", alias, interfaceType.String())
	}

	// getting outType pointer to check if implements the interface
	out := c.outType
	if out.Kind() != reflect.Ptr {
		out = reflect.New(out).Type()
	}

	// outType must implement the interface of the alias
	if !out.Implements(interfaceType) {
		log.Fatalf("[ERROR] '%s' not implements alias %s of type %s", out.String(), alias, interfaceType.String())
	}

	// Check if interfaceTypePtr is already registered for the instance
	if _, ok := c.interfaces[interfaceType]; ok {
		log.Fatalf("[ERROR] Alias '%s' of type '%s' already exists", alias, interfaceType.String())
	}

	// Check if alias is already registered for the instance
	if _, ok := c.aliases[alias]; ok {
		log.Fatalf("[ERROR] Alias '%s' already exists", alias)
	}

	al := transformAlias(alias)

	// Adding the alias
	c.aliases[al] = struct{}{}

	// Adding the new Interface for the instance
	c.interfaces[interfaceType] = struct{}{}

	return c
}

// AndInjections returns the instance to set the injections
func (c *container) AndInjections() InstanceInjection {
	return c
}

// WithInjection adds an injection alias to the parameter of the constructor
func (c *container) WithInjection(injectionAlias string, injectionTypePtr *reflect.Type) InstanceInjection {
	// Check if constructor have parameters
	if c.paramInjections.Len() == 0 {
		log.Fatalln("[ERROR] No injection parameters defined")
	}

	instanceType := *injectionTypePtr
	iter := c.paramInjections.Iterator()
	found := false
	for iter.Next() {
		v, pos, _ := iter.GetCurrentV()
		// If the type is the same as the injection type, and the injection is not already set
		if v == instanceType && c.paramInjections.Exists(pos) {
			// Replace the injection key (previous int position) with the alias
			c.paramInjections.ReplaceKey(pos, transformAlias(injectionAlias))
			found = true
			return c
		}
	}

	// If not found, it means that the injection type is not a parameter of the constructor
	if !found {
		log.Fatalf("[ERROR] Injection '%s' of type '%s' not found on params", injectionAlias, instanceType.String())
	}

	return c
}

// completeInjectionByType set the injection parameters and return if the constructor is ready to be called
func (c *container) completeInjectionByType(val reflect.Value) bool {
	iter := c.paramInjections.Iterator()
	for iter.Next() {
		t, i, _ := iter.GetCurrentV()
		if val.Type().ConvertibleTo(t) {
			c.paramsCount--
			c.params[i] = val
		}
	}
	return c.paramsCount == 0
}

// completeInjectionByAlias set the injection parameters and return if the constructor is ready to be called
func (c *container) completeInjectionByAlias(alias string, val reflect.Value) bool {
	t, _ := c.paramInjections.Get(alias)
	if val.Type().AssignableTo(t) {
		c.paramsCount--
		c.params[c.paramInjections.IndexOf(alias)] = val
	}
	return c.paramsCount == 0
}

// checkConstructor checks the constructor format
func checkConstructor(t reflect.Type) {
	if t.Kind() != reflect.Func {
		log.Fatalf("container must be a function")
	}
	if t.NumOut() != 1 {
		log.Fatalf("container must have one output")
	}
}

// transformAlias standardizes the alias
func transformAlias(alias string) string {
	alias = strings.ToLower(alias)
	alias = strings.Trim(alias, "")
	alias = strings.Replace(alias, " ", "-", -1)
	if len(alias) == 0 {
		log.Fatalf("[ERROR] Alias must not be empty")
	}
	return alias
}
