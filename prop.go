package lightms

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
)

var propFilePath = "resources/properties.yml"
var props = make([]interface{}, 0)

// SetPropFilePath sets the path of the property file. Default is "resources/properties.yml"
func SetPropFilePath(path string) {
	propFilePath = path
}

// AddProperty adds a property to the list of properties.
func AddProperty(p interface{}) {
	if pType := reflect.TypeOf(p); pType.Kind() != reflect.Ptr || pType.Elem().Kind() != reflect.Struct {
		log.Fatalf("Property must be a pointer to a struct. Got %v", pType)
	}
	props = append(props, p)
}

// loadProperties loads the properties from the yml property file.
func loadProperties() {
	log.Println("Started properties loading from file:", propFilePath)
	filename, err := filepath.Abs(propFilePath)
	if err != nil {
		log.Fatalf("Error while getting '%s' file path: %s\n", propFilePath, err)
	}

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error while reading '%s' file: %s\n", propFilePath, err)
	}
	yamlEnv := os.ExpandEnv(string(yamlFile))

	for _, p := range props {
		err = yaml.Unmarshal([]byte(yamlEnv), p)
		if err != nil {
			log.Fatalf("Error while unmarshalling '%s' file: %s\n", propFilePath, err)
		}
	}
	log.Println("Finished properties loading from file:", propFilePath)
}
