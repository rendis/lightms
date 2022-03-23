package secondary

import (
	"github.com/rendis/lightms/example/infra/config/prop"
	"log"
	"sync"
)

var (
	instance *PersistencePortImpl
	once     sync.Once
)

func GetPersistencePortImplInstance(prop prop.DataBaseInfo) *PersistencePortImpl {
	once.Do(func() {
		instance = &PersistencePortImpl{prop}
	})
	return instance
}

type PersistencePortImpl struct {
	prop prop.DataBaseInfo
}

func (p *PersistencePortImpl) Save(msg string) error {
	log.Printf("saving message '%s' in postgres database named '%s'.\n", msg, p.prop.Name)
	return nil
}
