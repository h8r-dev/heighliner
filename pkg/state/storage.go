package state

// Storage is an abstract of localstorage
type Storage interface {
	Detect() error
	GetPath() string
	Init() error
	Clean() error
}
