package lightms

import (
	"log"
	"reflect"
)

func newConfReader(addInstance func(string, reflect.Value), addDependent func(string, []reflect.Type, reflect.Value)) *confReader {
	return &confReader{
		confReg:      make([]any, 0),
		propRdr:      &propReader{},
		addInstance:  addInstance,
		addDependent: addDependent,
	}
}

type confReader struct {
	confReg      []any
	propRdr      *propReader
	addInstance  func(string, reflect.Value)
	addDependent func(string, []reflect.Type, reflect.Value)
}

func (c *confReader) addConf(conf any) {
	c.confReg = append(c.confReg, conf)
}

func (c *confReader) readConfReg() {
	for _, conf := range c.confReg {
		vconf := reflect.ValueOf(conf)
		validateConf(vconf.Type())
		c.loadConfProps(vconf)
		c.registerSupplier(vconf)
	}
}

func (c *confReader) loadConfProps(conf reflect.Value) {
	elem := conf.Type().Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		validateProp(elem.String(), field)
		propInstance := reflect.New(field.Type.Elem()).Interface()
		c.propRdr.loadProp(propInstance)
		log.Printf("Loaded prop '%s'", field.Name)
		vprop := reflect.ValueOf(propInstance)
		c.addInstance(field.Name, vprop)
		conf.Elem().Field(i).Set(vprop)
	}
}

func (c *confReader) registerSupplier(conf reflect.Value) {
	tconf := conf.Type()
	for i := 0; i < tconf.NumMethod(); i++ {
		method := tconf.Method(i)
		validateSupplier(conf.Type().String(), method.Type)
		numIn := method.Type.NumIn()
		supplier := conf.Method(i)
		if numIn == 1 {
			c.addInstance(method.Name, supplier.Call(nil)[0])
			continue
		}

		var dependTypes []reflect.Type
		for j := 1; j < numIn; j++ {
			dependTypes = append(dependTypes, method.Type.In(j))
		}
		c.addDependent(method.Name, dependTypes, supplier)
	}
}

func validateSupplier(confName string, supplierType reflect.Type) {
	if supplierType.NumOut() != 1 {
		log.Fatalf("Method '%s' in conf '%s' has wrong number of returns, expected 1 got %d.", supplierType.Name(), confName, supplierType.NumOut())
	}

	// check not duplicate param types
	var set = make(map[reflect.Type]bool)
	for i := 1; i < supplierType.NumIn(); i++ {
		_, exist := set[supplierType.In(i)]
		if exist {
			log.Fatalf("Method '%s' in conf '%s' has duplicate parameter type '%s'", supplierType.Name(), confName, supplierType.In(i).String())
		}
		set[supplierType.In(i)] = true
	}
}

func validateConf(tConf reflect.Type) {
	if tConf.Kind() != reflect.Ptr {
		log.Fatalf("Conf %s is not a pointer", tConf.String())
	}
}

func validateProp(confName string, field reflect.StructField) {
	if field.Type.Kind() != reflect.Ptr {
		log.Fatalf("Field '%s' in conf '%s' is not a pointer", field.Name, confName)
	}
}
