package impl

import (
	"github.com/rendis/lightms/example/core/port"
	"github.com/rendis/lightms/example/core/usecase"
	"log"
)

func NewJohnDoeImpl(p port.PersistencePort) usecase.JohnDoeUseCase {
	return &JohnDoeImpl{p}
}

type JohnDoeImpl struct {
	persistencePort port.PersistencePort
}

func (t *JohnDoeImpl) Handle(msg string) error {
	log.Printf("Usecase 'JohnDoeUseCase' handling message '%s'.\n", msg)
	return t.persistencePort.Save(msg)
}
