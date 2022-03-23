package instance

import (
	"github.com/rendis/lightms/example/core/impl"
	"github.com/rendis/lightms/example/core/usecase"
)

// GetJohnDoeUseCase returns impl.JohnDoeImpl instance
func GetJohnDoeUseCase() usecase.JohnDoeUseCase {
	return impl.GetJohnDoeImplInstance(
		GetPersistencePort(),
	)
}
