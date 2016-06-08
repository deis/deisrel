package docker

// Client is the interface to interact with the docker daemon.
type Client interface {
	Push(*Image) error
	Pull(*Image) error
	Retag(*Image, *Image) error
}
