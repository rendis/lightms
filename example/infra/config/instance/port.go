package instance

import (
	"github.com/rendis/lightms/example/core/port"
	"github.com/rendis/lightms/example/infra/config/prop"
	"github.com/rendis/lightms/example/infra/secondary"
)

// GetGetPersistencePort returns secondary.PersistencePortImpl instance
func GetGetPersistencePort() port.PersistencePort {
	return secondary.GetPersistencePortImplInstance(
		prop.GetDatabaseProp().Postgresql,
	)
}
