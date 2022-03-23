package port

type PersistencePort interface {
	Save(msg string) error
}
