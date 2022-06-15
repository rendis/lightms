package secondary

import (
	"github.com/rendis/lightms/example/core/port"
	"github.com/rendis/lightms/example/infra/config/prop"
	"log"
)

func NewPersistencePort(dbProp prop.DataBaseInfo) port.PersistencePort {
	return &PersistencePortImpl{dbProp}
}

type PersistencePortImpl struct {
	prop prop.DataBaseInfo
}

func (p *PersistencePortImpl) Save(msg string) error {
	log.Printf("Persistence port saving message '%s' in postgres database named '%s'.\n", msg, p.prop.Name)
	return nil
}
