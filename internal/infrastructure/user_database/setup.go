package user_database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectDatabase() {
	// @TODO - Implement env variables

	username := "root"
	password := "Oakbymarwin2545"
	// hostname := "rentless-product.chayaw1xzjuj.ap-southeast-1.rds.amazonaws.com"
	hostname := "localhost"
	port := "3306"
	dbname := "rentless"

	dsn := username + ":" + password + "@tcp(" + hostname + ":" + port + ")/" + dbname
	// dsn := fmt.Printf("%s:%s@tcp(%s:%s)/%s", username, password, hostname, port, dbname)
	fmt.Println(dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("error validating sql.Open arguments")
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("error verifying connection with %s", err)
		panic(err.Error())
	}

	fmt.Println("Successfully connected to MySQL database")

	defer db.Close()
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}
