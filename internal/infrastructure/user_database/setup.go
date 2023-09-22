package user_database

import (
	"database/sql"
	"fmt"

	models "rentless-services/internal/infrastructure/user_database/models"

	_ "github.com/go-sql-driver/mysql"
)

type UserDB struct {
	DB *sql.DB
}

func ConnectDatabase() *UserDB {
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

	// Connect Database
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
	return &UserDB{
		DB: db,
	}
}

var db *sql.DB

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// Get One User done
func (userDB *UserDB) GetOneUser(id string) (models.User, error) {
	query := "SELECT id, firstname, lastname, age, created_at, updated_at FROM users WHERE id = ?"

	var user models.User
	err := userDB.DB.QueryRow(query, id).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Age, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// Get All User done
func (userDB *UserDB) GetAllUsers() ([]models.User, error) {
	query := "SELECT id, firstname, lastname, age, created_at, updated_at FROM users"

	// Execute the query
	rows, err := userDB.DB.Query(query)
	if err != nil {
		return []models.User{}, err
	}

	// Initialize a slice to hold the result
	var users []models.User

	// Loop through the result set and scan each row into a User struct
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Age, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Create User done
func (userDB *UserDB) CreateUser(user models.User) (models.User, error) {
	query := "INSERT INTO users (firstname, lastname, age, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())"

	// Execute the query
	_, err := userDB.DB.Exec(query, user.Firstname, user.Lastname, user.Age)
	if err != nil {
		return models.User{}, err
	}

	query = "SELECT id, firstname, lastname, age FROM users ORDER BY id DESC LIMIT 1"
	var new_user models.User
	err = userDB.DB.QueryRow(query).Scan(&user.ID, &new_user.Firstname, &new_user.Lastname, &new_user.Age)
	if err != nil {
		return models.User{}, err
	}
	return new_user, nil
}

// Update User done
func (userDB *UserDB) UpdateUser(user models.User) (models.User, error) {
	query := "UPDATE users SET firstname = ?, lastname = ?, age = ?, updated_at = NOW() WHERE id = ?"

	_, err := userDB.DB.Exec(query, user.Firstname, user.Lastname, user.Age, user.ID)
	if err != nil {
		return models.User{}, err
	}

	query = "SELECT id, firstname, lastname, age FROM users WHERE id = ?"
	var new_user models.User
	err = userDB.DB.QueryRow(query, user.ID).Scan(&new_user.ID, &new_user.Firstname, &new_user.Lastname, &new_user.Age)
	if err != nil {
		return models.User{}, err
	}
	return new_user, nil
}

// Delete User done
func (userDB *UserDB) DeleteUser(id string) error {
	query := "DELETE FROM users WHERE id = ?"

	_, err := userDB.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// Is User Exist
func (userDB *UserDB) userExists(id string) bool {
	query := "SELECT id FROM users WHERE id = ?"
	// var userID uint
	// err := userDB.DB.QueryRow(query, id).Scan(&userID)
	// return err != sql.ErrNoRows
	var count int
	fmt.Println("111111111")
	err := userDB.DB.QueryRow(query, id).Scan(&count)
	if err != nil {
		fmt.Println("Scan Error")
		return false
	}
	fmt.Println("22222222222")
	return count > 0
}
