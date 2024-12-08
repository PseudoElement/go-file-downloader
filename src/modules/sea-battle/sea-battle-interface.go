package seabattle

type PlayerSocket interface {
	Connect() error
	Disconnect() error
	Broadcast()
}
