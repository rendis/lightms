package lightms

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	propPathEnv     = "APP_CONFIG_FILE"
	defaultPropPath = "resources/properties.yml"
)

var (
	propFilePath = os.Getenv(propPathEnv)
)

// SetPropFilePath sets the path of the property file. Default is "resources/properties.yml"
func SetPropFilePath(propPath string) {
	if propPath == "" {
		log.Fatalf("Yml file  parameter is empty.")
	}
	propFilePath = propPath
}

type PropDefault interface {
	SetDefault()
}

type propReader struct {
	loadPropsOnce sync.Once
	propArr       []byte
}

func (p *propReader) loadProp(prop any) {
	p.readYml()
	env := []byte(os.ExpandEnv(string(p.propArr)))
	err := yaml.Unmarshal(env, prop)
	p.runDefault(prop)
	if err != nil {
		log.Fatalf("Error parsing yml file '%s' to struct '%v'. %s", propFilePath, prop, err)
	}
}

func (p *propReader) runDefault(prop any) {
	if d, ok := prop.(PropDefault); ok {
		d.SetDefault()
	}
}

func (p *propReader) readYml() {
	p.loadPropsOnce.Do(func() {
		if propFilePath == "" {
			log.Printf("%s is not set, using default supplier: %s\n", propPathEnv, defaultPropPath)
			propFilePath = defaultPropPath
		}

		filename, err := filepath.Abs(propFilePath)
		if err != nil {
			log.Fatalf("Error getting yml file '%s'. %s", filename, err)
		}

		b, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalf("Error reading yml file '%s'. %s", filename, err)
		}
		p.propArr = b
	})
}
