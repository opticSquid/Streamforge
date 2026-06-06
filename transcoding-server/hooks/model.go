package hooks

type StreamManager interface {
	Start(name string) error
	Stop(name string)
}

type Handler struct {
	manager StreamManager
}
