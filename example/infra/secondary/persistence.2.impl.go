package secondary

import (
	"github.com/rendis/lightms/v3/example/infra/config/prop"
	"log"
)

func NewPersistence2Port(dbProp *prop.DatabaseProp) *PersistencePort2Impl {
	return &PersistencePort2Impl{dbProp.Postgresql}
}

type PersistencePort2Impl struct {
	prop *prop.DataBaseInfo
}

func (p *PersistencePort2Impl) Save(msg string) error {
	log.Printf("Persistence port saving message '%s' in postgres database named '%s'.\n", msg, p.prop.Name)
	return nil
}
