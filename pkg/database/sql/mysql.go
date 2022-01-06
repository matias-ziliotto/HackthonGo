package sql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	MySqlDB *sql.DB
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

	MySqlDB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	if err = MySqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("DB ready")
}
