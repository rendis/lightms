package impl

import (
	"github.com/rendis/lightms/example/core/port"
	"sync"
)

var (
	instance *JohnDoeImpl
	once     sync.Once
)

func GetJohnDoeImplInstance(p port.PersistencePort) *JohnDoeImpl {
	once.Do(func() {
		instance = &JohnDoeImpl{p}
	})
	return instance
}

type JohnDoeImpl struct {
	persistencePort port.PersistencePort
}

func (t *JohnDoeImpl) Handle(msg string) error {
	return t.persistencePort.Save(msg)
}
