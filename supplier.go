package lightms

import (
	"log"
	"reflect"
)

type supplierInfo struct {
	name      string
	inCount   int
	supplier  reflect.Value
	inputType []reflect.Type
	input     []reflect.Value
}

func (g *supplierInfo) addDependency(val reflect.Value) int {
	var index = g.indexOfType(val.Type())
	if index == -1 {
		log.Fatalf("Wrong dependency injection for %s, %s", g.name, val.Type())
	}
	g.inputType[index] = reflect.TypeOf(nil)
	g.input[index] = val
	g.inCount--
	return g.inCount
}

func (g *supplierInfo) indexOfType(t reflect.Type) int {
	for i, inputType := range g.inputType {
		if inputType == t {
			return i
		}
	}
	return -1
}

func (g *supplierInfo) call() reflect.Value {
	return g.supplier.Call(g.input)[0]
}
