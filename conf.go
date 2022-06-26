package lightms

import (
	"log"
	"reflect"
)

var propsConf []any
var reader = newPropReader()

// AddConf registers a conf struct to be loaded
func AddConf[T any]() {
	t := new(T)
	typ := reflect.TypeOf(t).Elem()
	if typ.Kind() != reflect.Struct {
		log.Fatalf("AddConf: conf must be a struct, got '%s'", typ.Kind())
	}
	propsConf = append(propsConf, t)
}

// processInjectionsByReceivers analyze receivers to populate dByAlias and dByTypes
func processInjectionsByReceivers() {
	for _, conf := range propsConf {
		loadConfMeth(reflect.ValueOf(conf))
	}
}

// resolveProps resolves all props fields
func resolveProps() {
	for _, conf := range propsConf {
		loadConfProps(reflect.ValueOf(conf))
	}
}

// loadConfMeth calls receiver method
func loadConfMeth(vconf reflect.Value) {
	tconf := vconf.Type()
	for i := 0; i < tconf.NumMethod(); i++ {
		method := tconf.Method(i)
		validateConfMeth(vconf.Type().String(), method.Name, method.Type)
		vconf.Method(i).Call(nil)
	}
}

// validateConfMeth checks if the receiver method
func validateConfMeth(confName, methName string, methType reflect.Type) {
	if methType.NumOut() != 0 {
		log.Fatalf("Method '%s' in conf '%s' has wrong number of returns, expected 0 got %d.", methName, confName, methType.NumOut())
	}

	if methType.NumIn() != 1 {
		log.Fatalf("Method '%s' in conf '%s' has wrong number of inputs, expected 0 got %d.", methName, confName, methType.NumIn())
	}
}

// loadConfProps loads conf props
func loadConfProps(conf reflect.Value) {
	elem := conf.Type().Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		log.Printf("Loading prop %s", field.Name)
		propInstance := reflect.New(field.Type.Elem()).Interface()
		reader.loadProp(propInstance)
		vprop := reflect.ValueOf(propInstance)
		storePropResolved(vprop)
		conf.Elem().Field(i).Set(vprop)
	}
}
