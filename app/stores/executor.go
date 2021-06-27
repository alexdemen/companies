package stores

type Executor interface {
	Close(err error)
	AddBuilding() error
}
