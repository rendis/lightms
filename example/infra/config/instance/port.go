package instance

import (
	"github.com/rendis/lightms/example/core/port"
	"github.com/rendis/lightms/example/infra/config/prop"
	"github.com/rendis/lightms/example/infra/secondary"
)

// GetPersistencePort returns secondary.PersistencePortImpl instance
func GetPersistencePort() port.PersistencePort {
	return secondary.GetPersistencePortImplInstance(
		prop.GetDatabaseProp().Postgresql,
	)
}
