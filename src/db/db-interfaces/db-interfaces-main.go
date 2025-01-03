package db_interfaces

type Database[T any] interface {
	Connect() error
	Disconnect() error
	Conn() T
}

type TableCreator interface {
	CreateTables() error
}
