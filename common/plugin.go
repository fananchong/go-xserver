package common

// Plugin : The interface that conforms to the go-xserver call
type Plugin interface {
	Init() bool
	Start() bool
	Close()
}
