package review_database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	models "rentless-services/internal/infrastructure/review_database/models"

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
func (reviewDB *ReviewDB) GetOneRecord(review_id string) (*models.Review, error) {
	query := "SELECT id, created_at, author_id, product_id, rate, comment FROM reviews WHERE id = ?"
	var review models.Review
	err := reviewDB.DB.QueryRow(query, review_id).Scan(&review.ID, &review.CreatedAt, &review.AuthorID, &review.ProductID, &review.Rate, &review.Comment)
	if err != nil {
		return nil, err
	}
	return &review, nil
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

// Get Average of rate and number of reviews by product_id
func (reviewDB *ReviewDB) GetAvgRateAndCountByProductID(product_id string) (*models.AvgAndCount, error) {
	query := "SELECT COALESCE(AVG(rate), 0) AS average_rate, COUNT(*) AS count_reviews FROM reviews WHERE product_id = ?"
	var avgAndCount models.AvgAndCount
	err := reviewDB.DB.QueryRow(query, product_id).Scan(&avgAndCount.AvgRate, &avgAndCount.NumberReview)
	if err != nil {
		return nil, err
	}
	return &avgAndCount, nil
}

// Get Average of rate and number of reviews
func (reviewDB *ReviewDB) GetAvgRateAndCountAll() (*models.AvgAndCount, error) {
	query := "SELECT COALESCE(AVG(rate), 0) AS average_rate, COUNT(*) AS count_reviews FROM reviews"
	var avgAndCount models.AvgAndCount
	err := reviewDB.DB.QueryRow(query).Scan(&avgAndCount.AvgRate, &avgAndCount.NumberReview)
	if err != nil {
		return nil, err
	}
	return &avgAndCount, nil
}

// Get reviews Exist
func (reviewDB *ReviewDB) ReviewExists(id string) bool {
	query := "SELECT id FROM reviews WHERE id = ?"
	var userID uint
	err := reviewDB.DB.QueryRow(query, id).Scan(&userID)
	return err != sql.ErrNoRows
}

// Create Records
func (reviewDB *ReviewDB) CreateRecord(author_id, product_id uint, review models.Review) error {
	query := "INSERT INTO reviews (created_at, author_id, product_id, rate, comment) VALUES (NOW(), ?, ?, ?, ?)"
	_, err := reviewDB.DB.Exec(query, author_id, product_id, review.Rate, review.Comment)
	if err != nil {
		return err
	}
	return nil
}
