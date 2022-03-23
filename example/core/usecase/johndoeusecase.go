package usecase

type JohnDoeUseCase interface {
	Handle(msg string) error
}
