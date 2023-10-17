package review_database

import (
	"database/sql"
	"fmt"

	models "rentless-services/internal/infrastructure/review_database/models"

	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type ReviewDB struct {
	DB *sql.DB
}

func init() {

	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
		// panic(err.Error())
	}
}

func ConnectDatabase() *ReviewDB {
	// dsn := "root:Oakbymarwin2545@tcp(localhost:3306)/rentless?parseTime=true"
	dsn := os.Getenv("REVIEW_DB_URL")
	fmt.Printf("URL: %s \n", dsn)

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

	return &ReviewDB{
		DB: db,
	}
}

// Get One Record By Review ID
func (reviewDB *ReviewDB) GetOneRecord(id string) ([]models.Review, error) {
	query := "SELECT id, created_at, author_id, product_id, rate, comment FROM reviews WHERE id = ?"
	var review models.Review
	err := reviewDB.DB.QueryRow(query, id).Scan(&review.ID, &review.CreatedAt, &review.AuthorID, &review.ProductID, &review.Rate, &review.Comment)
	if err != nil {
		return nil, err
	}
	var reviews []models.Review
	reviews = append(reviews, review)
	return reviews, nil

}

// Get Records By Product ID
func (reviewDB *ReviewDB) GetRecordsByProductID(product_id string) ([]models.Review, error) {
	query := "SELECT id, created_at, author_id, product_id, rate, comment FROM reviews WHERE product_id = ?"
	rows, err := reviewDB.DB.Query(query, product_id)
	if err != nil {
		return nil, err
	}
	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		err := rows.Scan(&review.ID, &review.CreatedAt, &review.AuthorID, &review.ProductID, &review.Rate, &review.Comment)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

// Get All Records
func (reviewDB *ReviewDB) GetAllRecords() ([]models.Review, error) {
	query := "SELECT id, created_at, author_id, product_id, rate, comment FROM reviews"
	rows, err := reviewDB.DB.Query(query)
	if err != nil {
		return nil, err
	}
	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		err := rows.Scan(&review.ID, &review.CreatedAt, &review.AuthorID, &review.ProductID, &review.Rate, &review.Comment)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

// Create Records
func (reviewDB *ReviewDB) CreateRecord(review models.Review) error {
	query := "INSERT INTO reviews (created_at, author_id, product_id, rate, comment) VALUES (NOW(), ?, ?, ?, ?)"
	_, err := reviewDB.DB.Exec(query, review.AuthorID, review.ProductID, review.Rate, review.Comment)
	if err != nil {
		return err
	}
	return nil
}
