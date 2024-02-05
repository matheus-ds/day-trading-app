package store

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgresClient struct {
	conn *sqlx.DB
}

func (pgc *PostgresClient) Close() {
	pgc.conn.Close()
}

func NewPostgresClient(dbURL string) *PostgresClient {
	conn, err := sqlx.Connect("pgx", dbURL)

	if err != nil {
		log.Println(err)
	}

	pingErr := conn.Ping()

	if pingErr != nil {
		log.Println(pingErr)
	}

	log.Println("Connected!")

	return &PostgresClient{conn: conn}
}
