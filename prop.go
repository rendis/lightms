package lightms

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	propPathEnv     = "APP_CONFIG_FILE"
	defaultPropPath = "resources/properties.yml"
)

var (
	propFilePath = os.Getenv(propPathEnv)
)

type ConfigFileType int

const (
	Yml ConfigFileType = iota
	Json
)

var ErrInvalidConfigFileType = errors.New("invalid config file type. only yml and json are supported")

// SetPropFilePath sets the path of the property file. Default is "resources/properties.yml"
// property file can be a yml file or a json file
func SetPropFilePath(propPath string) {
	propPath = strings.Trim(propPath, " ")
	if propPath == "" {
		log.Fatalf("Yml file  parameter is empty.")
	}

	if _, err := getConfigFileType(propPath); err != nil {
		log.Fatalf("Error getting config file type. %s", err)
	}

	propFilePath = propPath
}

// getConfigFileType returns the config file type
// property file can be a yml file or a json file
func getConfigFileType(path string) (ConfigFileType, error) {
	ext := filepath.Ext(path)
	switch ext {
	case ".yml":
		return Yml, nil
	case ".json":
		return Json, nil
	default:
		return 0, ErrInvalidConfigFileType
	}
}

func getConfPathAndType() (string, ConfigFileType, error) {
	if propFilePath == "" {
		log.Printf("%s is not set, using default supplier: %s\n", propPathEnv, defaultPropPath)
		propFilePath = defaultPropPath
	}

	filename, err := filepath.Abs(propFilePath)
	if err != nil {
		return "", 0, err
	}

	fileType, err := getConfigFileType(filename)
	if err != nil {
		return "", 0, err
	}

	return filename, fileType, nil
}

type PropDefault interface {
	SetDefault()
}

func newPropReader() *propReader {
	p := &propReader{}
	return p
}

type propReader struct {
	propArr        []byte
	configFileType ConfigFileType
	loadPropsOnce  sync.Once
}

func (p *propReader) loadProp(prop any) {
	p.readConfigFile()
	switch p.configFileType {
	case Yml:
		p.readFromYml(prop)
	case Json:
		p.readFromJson(prop)
	}
	p.runDefault(prop)
}

func (p *propReader) runDefault(prop any) {
	if d, ok := prop.(PropDefault); ok {
		d.SetDefault()
	}
}

func (p *propReader) readConfigFile() {
	p.loadPropsOnce.Do(func() {
		fileName, typ, err := getConfPathAndType()
		if err != nil {
			log.Fatalf("Error getting config file type. %s", err)
		}

		b, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatalf("Error reading config file. %s", err)
		}
		propWithEnv := []byte(os.ExpandEnv(string(b)))
		p.propArr = propWithEnv

		// propWithEnv to base64
		b64 := base64.StdEncoding.EncodeToString(propWithEnv)
		fmt.Printf("Property file content B64: %s\n\n", b64)
		fmt.Printf("Property file content:\n%s", string(propWithEnv))

		p.configFileType = typ
	})
}

func (p *propReader) readFromYml(prop any) {
	if err := yaml.Unmarshal(p.propArr, prop); err != nil {
		log.Fatalf("Error parsing yml file '%s' to struct '%v'. %s", propFilePath, prop, err)
	}
}

func (p *propReader) readFromJson(prop any) {
	if err := json.Unmarshal(p.propArr, prop); err != nil {
		log.Fatalf("Error parsing json file '%s' to struct '%v'. %s", propFilePath, prop, err)
	}
}
