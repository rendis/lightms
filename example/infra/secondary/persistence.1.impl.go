package secondary

import (
	"github.com/rendis/lightms/v3/example/infra/config/prop"
	"log"
)

func NewPersistence1Port(dbProp *prop.DatabaseProp) *PersistencePort1Impl {
	return &PersistencePort1Impl{dbProp.Postgresql}
}

type PersistencePort1Impl struct {
	prop *prop.DataBaseInfo
}

func (p *PersistencePort1Impl) Save(msg string) error {
	log.Printf("Persistence port saving message '%s' in postgres database named '%s'.\n", msg, p.prop.Name)
	return nil
}
