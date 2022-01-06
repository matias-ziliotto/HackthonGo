package sql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/DATA-DOG/go-txdb"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../../cmd/server/.env")

	if err != nil {
		log.Fatal(err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbServer := os.Getenv("DB_HOST")
	dbHost := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbServer, dbHost, dbName)

	txdb.Register("txdb", "mysql", connectionString)
}

func InitTxSqlDb() (*sql.DB, error) {
	db, err := sql.Open("txdb", uuid.New().String())
	if err == nil {
		return db, db.Ping()
	}

	return db, err
}
