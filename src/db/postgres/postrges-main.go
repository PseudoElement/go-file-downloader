package postgres

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	db_interfaces "github.com/pseudoelement/go-file-downloader/src/db/db-interfaces"
)

type PostgresDB struct {
	user     string
	dbName   string
	password string
	host     string
	conn     *sqlx.DB
}

func New() db_interfaces.Database[*sqlx.DB] {
	return &PostgresDB{
		user:     "postgres",
		dbName:   "downloader_db",
		password: "password",
		host:     "localhost",
	}
}

func (db *PostgresDB) Connect() error {
	connData := fmt.Sprintf(
		"user=%s dbname=%s sslmode=disable password=%s host=%s",
		db.user,
		db.dbName,
		db.password,
		db.host,
	)

	conn, err := sqlx.Connect("postgres", connData)
	if err != nil {
		panic(err)
	}

	db.conn = conn

	log.Println("PostgresDB successfully connected!")

	return nil
}

func (db *PostgresDB) Disconnect() error {
	err := db.conn.Close()
	return err
}

func (db *PostgresDB) Conn() *sqlx.DB {
	return db.conn
}

var _ db_interfaces.Database[*sqlx.DB] = (*PostgresDB)(nil)
