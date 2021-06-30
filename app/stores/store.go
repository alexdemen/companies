package stores

type Store interface {
	GetExecutor() (Executor, error)
}
