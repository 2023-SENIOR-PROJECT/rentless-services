package db

import (
	"authservice/models"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type UserDB struct {
	DB *sql.DB
}

func ConnectDatabase() *UserDB {
	// @TODO - Implement env variables

	username := "root"
	password := "password"
	// hostname := "rentless-product.chayaw1xzjuj.ap-southeast-1.rds.amazonaws.com"
	hostname := "localhost"
	port := "3307"
	dbname := "rentless"

	dsn := username + ":" + password + "@tcp(" + hostname + ":" + port + ")/" + dbname + "?parseTime=true"
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

	return &UserDB{
		DB: db,
	}
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

	query = "SELECT id, firstname, lastname, age, created_at, updated_at FROM users ORDER BY id DESC LIMIT 1"
	var new_user models.User
	err = userDB.DB.QueryRow(query).Scan(&new_user.ID, &new_user.Firstname, &new_user.Lastname, &new_user.Age, &new_user.CreatedAt, &new_user.UpdatedAt)
	if err != nil {
		return models.User{}, err
	}
	return new_user, nil
}

// Update User done
func (userDB *UserDB) UpdateUser(id string, user models.User) (models.User, error) {
	query := "UPDATE users SET firstname = ?, lastname = ?, age = ?, updated_at = NOW() WHERE id = ?"

	_, err := userDB.DB.Exec(query, user.Firstname, user.Lastname, user.Age, id)
	if err != nil {
		return models.User{}, err
	}
	fmt.Println(id)
	query = "SELECT id, firstname, lastname, age FROM users WHERE id = ?"
	var new_user models.User
	err = userDB.DB.QueryRow(query, id).Scan(&new_user.ID, &new_user.Firstname, &new_user.Lastname, &new_user.Age)
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

// Is User Exist Rewrite
func (userDB *UserDB) UserExists(id string) bool {
	query := "SELECT id FROM users WHERE id = ?"
	var userID uint
	err := userDB.DB.QueryRow(query, id).Scan(&userID)
	return err != sql.ErrNoRows
}

// Auth service
func (userDB *UserDB) CreateNewUser(user models.UserAuthSturct, user_id uint) error {

	query := "INSERT INTO auths (email, pwd, token, user_id) VALUES (?, ?, ?, ?)"
	_, err := userDB.DB.Exec(query, user.Email, user.Pwd, user.Token, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (userDB *UserDB) FindUserAccount(req models.LoginRequest) (models.UserAuthSturct, error) {
	query := "SELECT email, pwd, token, user_id FROM auths WHERE email = ? AND pwd = ?"
	var user models.UserAuthSturct
	err := userDB.DB.QueryRow(query, req.Email, req.Pwd).Scan(&user.Email, &user.Pwd, &user.Token, &user.User_id)
	if err != nil {
		return models.UserAuthSturct{}, err
	}
	return user, nil
}
