package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

const (
	driverName = "postgres"
	dataSourceName = "postgresql://root:123@localhost:5432/simple_bank?sslmode=disable"
)


func TestMain(m *testing.M){
	conn, err := sql.Open(driverName, dataSourceName)

	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}