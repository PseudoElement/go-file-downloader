package seabattle

type PlayerSocket interface {
	Connect() error
	Disconnect(errMsg *ErrorForDB) error
	Broadcast()
}
