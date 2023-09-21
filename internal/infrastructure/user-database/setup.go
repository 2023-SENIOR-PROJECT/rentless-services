// >> Unfinished code << 

package user_database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDatabase() {
	// @TODO - Implement env variables
	db, err := sql.Open(("mysql"), "rentless-product.chayaw1xzjuj.ap-southeast-1.rds.amazonaws.com")
	if err != nil {
		fmt.Println("error validating sql.Open arguments")
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("error verifying connection with db.Ping")
		panic(err.Error())
	}

	fmt.Println("Successfully connected to MySQL database")
}
