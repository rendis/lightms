package lightms

import (
	"fmt"
	"github.com/rendis/orderedmap"
	"log"
	"reflect"
	"strings"
)

var ppcounter = 0

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
	return registerInstance(f)
}

// AddInstanceWithAlias adds an instance to the container with an alias
func AddInstanceWithAlias(alias string, f any) InstanceAlias {
	return registerInstance(f, alias)
}

func registerInstance(f any, alias ...string) InstanceAlias {
	v := reflect.ValueOf(f)
	t := v.Type()
	checkConstructor(t)
	out := t.Out(0)

	// Getting the alias
	var al = ""
	switch {
	case len(alias) == 1:
		al = transformAlias(alias[0])
	case out.Kind() == reflect.Ptr:
		al = transformAlias(out.Elem().String())
	default:
		al = transformAlias(out.String())
	}

	// check if is a Primary Process
	isPP := false
	if out.Implements(pptype) {
		isPP = true
	}

	c := &container{
		constructorFunc: v,
		outType:         out,
		isPP:            isPP,
		aliases:         []string{al},
		paramsCount:     t.NumIn(),
		params:          make([]reflect.Value, t.NumIn()),
		paramInjections: orderedmap.New[reflect.Type](),
		aliasTypes:      make(map[reflect.Type]struct{}),
	}

	// Registering the params to inject
	for i := 0; i < t.NumIn(); i++ {
		c.paramInjections.Set(i, t.In(i))
	}

	// If is a Primary Process and alias is not defined, it will be the default alias
	if isPP && len(alias) == 0 {
		ppcounter++
		al = fmt.Sprintf("%s-%d", al, ppcounter)
	}

	// First alias is the same instance
	c.aliasTypes[out] = struct{}{}

	// Registering the instance
	addDependency(c)
	addAliasResolver(al, c)
	addTypeResolver(out, c)
	return c
}

type container struct {
	constructorFunc reflect.Value                       // constructor function
	outType         reflect.Type                        // output type
	isPP            bool                                // is Primary Process
	aliases         []string                            // aliases of the instance
	aliasTypes      map[reflect.Type]struct{}           // to avoid alias type repetition
	params          []reflect.Value                     // input parameters to pass to the constructor
	paramsCount     int                                 // number of remaining parameters to pass to the constructor
	paramInjections orderedmap.OrderedMap[reflect.Type] // ordered map used to know how param are injected by alias or type
}

// WithAlias adds an alias with a type to the instance
func (r *container) WithAlias(alias string, ptrInstanceType *reflect.Type) InstanceAlias {
	instanceType := *ptrInstanceType

	// instanceType must be an interface
	if instanceType.Kind() != reflect.Interface {
		log.Fatalf("[ERROR] Alias must be of type interface. Alias: %s, Type: %s", alias, instanceType.String())
	}

	// Getting outType pointer
	out := r.outType
	if out.Kind() != reflect.Ptr {
		out = reflect.New(out).Type()
	}

	// outType must implement the interface of the alias
	if !out.Implements(instanceType) {
		log.Fatalf("[ERROR] '%s' not implements alias %s of type %s", out.String(), alias, instanceType.String())
	}

	// Check if ptrInstanceType is already registered for the instance
	if _, ok := r.aliasTypes[instanceType]; ok {
		log.Fatalf("[ERROR] Alias '%s' of type '%s' already exists", alias, instanceType.String())
	}

	al := transformAlias(alias)

	// Registering the instance
	addAliasResolver(al, r)
	addTypeResolver(instanceType, r)

	// Adding the alias
	r.aliases = append(r.aliases, al)

	// Adding the new type to the instance
	r.aliasTypes[instanceType] = struct{}{}

	return r
}

// AndInjections returns the instance to set the injections
func (r *container) AndInjections() InstanceInjection {
	return r
}

// WithInjection adds an injection alias to the parameter of the constructor
func (r *container) WithInjection(injectionAlias string, ptrInstanceType *reflect.Type) InstanceInjection {
	// Check if constructor have parameters
	if r.paramInjections.Len() == 0 {
		log.Fatalln("[ERROR] No injection parameters defined")
	}

	instanceType := *ptrInstanceType
	iter := r.paramInjections.Iterator()
	found := false
	for iter.Next() {
		v, pos, _ := iter.GetCurrentV()
		// If the type is the same as the injection type, and the injection is not already set
		if v == instanceType && r.paramInjections.Exists(pos) {
			// Replace the injection key (previous int position) with the alias
			r.paramInjections.ReplaceKey(pos, transformAlias(injectionAlias))
			found = true
			return r
		}
	}

	// If not found, it means that the injection type is not a parameter of the constructor
	if !found {
		log.Fatalf("[ERROR] Injection '%s' of type '%s' not found on params", injectionAlias, instanceType.String())
	}

	return r
}

// completeInjectionByType set the injection parameters and return if the constructor is ready to be called
func (r *container) completeInjectionByType(val reflect.Value) bool {
	iter := r.paramInjections.Iterator()
	for iter.Next() {
		t, i, _ := iter.GetCurrentV()
		if val.Type() == t {
			r.paramsCount--
			r.params[i] = val
		}
	}
	return r.paramsCount == 0
}

// completeInjectionByAlias set the injection parameters and return if the constructor is ready to be called
func (r *container) completeInjectionByAlias(alias string, val reflect.Value) bool {
	t, _ := r.paramInjections.Get(alias)
	if val.Type().AssignableTo(t) {
		r.paramsCount--
		r.params[r.paramInjections.IndexOf(alias)] = val
	}
	return r.paramsCount == 0
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
