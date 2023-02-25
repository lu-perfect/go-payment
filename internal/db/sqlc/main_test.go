package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

// TODO: move to configuration

const (
	DBDriver = "postgres"
	DBSource = "postgresql://root:secret@localhost:5432/gobank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	testDB, err := sql.Open(DBDriver, DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
