package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	db_interfaces "github.com/pseudoelement/go-file-downloader/src/db/db-interfaces"
)

type PostgresDB struct {
	user     string
	dbName   string
	password string
	host     string
	port     int32
	conn     *sql.DB
}

func New() db_interfaces.Database[*sql.DB] {
	return &PostgresDB{
		// user name of desktop, with docker it's `postgres`
		user:     "paveldavidovich",
		dbName:   "postgres",
		password: "postgres",
		host:     "localhost",
		// host:     "postgres", with docker
		port: 5432,
	}
}

func (db *PostgresDB) Connect() error {
	connData := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db.host,
		db.port,
		db.user,
		db.password,
		db.dbName,
	)

	conn, err := sql.Open("postgres", connData)
	if err != nil {
		panic(err)
	}

	err = conn.Ping()
	if err != nil {
		panic(err)
	}

	db.conn = conn

	db.createTestTableIfNotExists()

	log.Println("PostgresDB successfully connected!")

	return nil
}

func (db *PostgresDB) Disconnect() error {
	err := db.conn.Close()
	return err
}

func (db *PostgresDB) Conn() *sql.DB {
	return db.conn
}

func (db *PostgresDB) createTestTableIfNotExists() {
	_, err := db.conn.Exec("CREATE TABLE IF NOT EXISTS test_table(last_name varchar(255), first_name varchar(255));")
	if err != nil {
		msg := fmt.Sprintf("Error creating test_table - %v", err)
		panic(msg)
	}
	log.Println("Table created!")
}

var _ db_interfaces.Database[*sql.DB] = (*PostgresDB)(nil)
