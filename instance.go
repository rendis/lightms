package lightms

import (
	"log"
	"reflect"
	"sync"
)

func newInstanceRegistry() *instanceRegistry {
	return &instanceRegistry{
		dependentSuppliers: make(map[string]*supplierInfo),
		namedRegistry:      make(map[string]reflect.Value),
		typedRegistry:      make(map[reflect.Type][]typedInstance),
		instanceNames:      make([]reflect.Type, 0),
	}
}

type InstanceProvider struct {
	GetByName func(string) any
}

type instanceRegistry struct {
	registryMtx        sync.Mutex
	dependentSuppliers map[string]*supplierInfo
	namedRegistry      map[string]reflect.Value
	typedRegistry      map[reflect.Type][]typedInstance
	instanceNames      []reflect.Type
}

func (i *instanceRegistry) getInstanceByName(name string) any {
	if instance, ok := i.namedRegistry[name]; ok {
		return instance.Interface()
	}
	return nil
}

func (i *instanceRegistry) addInstance(name string, value reflect.Value) {
	i.registryMtx.Lock()
	defer i.registryMtx.Unlock()

	_, ok := i.namedRegistry[name]
	if ok {
		log.Fatalf("Instance '%s' already registered", name)
	}
	i.namedRegistry[name] = value
	i.typedRegistry[value.Type()] = append(i.typedRegistry[value.Type()], typedInstance{name, value})
	i.instanceNames = append(i.instanceNames, value.Type())

	isPP := value.Type().Implements(reflect.TypeOf((*PrimaryProcess)(nil)).Elem())
	if isPP {
		addPrimary(name, func() PrimaryProcess {
			return value.Interface().(PrimaryProcess)
		})
	}
}

func (i *instanceRegistry) addDependent(instanceName string, dependTypes []reflect.Type, supplier reflect.Value) {
	if i.dependentSuppliers[instanceName] != nil {
		log.Fatalf("Instance '%s' already registered", instanceName)
	}

	for _, t := range dependTypes {
		if t.Implements(reflect.TypeOf((*PrimaryProcess)(nil)).Elem()) {
			log.Fatalf("Instances cannot depend on primary process (PrimaryProcess). Instance name: %s", instanceName)
		}
	}

	d := &supplierInfo{
		name:      instanceName,
		supplier:  supplier,
		inCount:   len(dependTypes),
		input:     make([]reflect.Value, len(dependTypes)),
		inputType: dependTypes,
	}
	i.dependentSuppliers[instanceName] = d
}

func (i *instanceRegistry) getDependsByType(t reflect.Type) []*supplierInfo {
	var result []*supplierInfo
	for _, s := range i.dependentSuppliers {
		if s.indexOfType(t) != -1 {
			result = append(result, s)
		}
	}
	return result
}

func (i *instanceRegistry) resolveDependents() {
	if len(i.dependentSuppliers) != 0 {
		for len(i.instanceNames) != 0 {
			ttype := i.instanceNames[0]
			i.instanceNames = i.instanceNames[1:]
			instances := i.typedRegistry[ttype]
			for _, d := range i.getDependsByType(ttype) {
				if d.addDependency(instances[0].value) == 0 {
					log.Printf("Instance '%s' is ready", d.name)
					i.addInstance(d.name, d.call())
					delete(i.dependentSuppliers, d.name)
				}
			}
		}
	}

	if len(i.dependentSuppliers) != 0 {
		log.Fatalf("Dependent suppliers not resolved: %v", i.dependentSuppliers)
	}
}

type typedInstance struct {
	name  string
	value reflect.Value
}
