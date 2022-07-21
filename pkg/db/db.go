package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//Initialize postgres database
func Init() *sql.DB {
	//load dotenv file
	err := godotenv.Load()
	checkErr(err)

	//database config
	host := "localhost"
	port := 5432
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASS")
	dbname := os.Getenv("DBNAME")

	// Get a database handle.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	checkErr(err)

	// Connect to database
	pingErr := db.Ping()
	checkErr(pingErr)

	fmt.Println("Connected to DB!")

	return db
}
