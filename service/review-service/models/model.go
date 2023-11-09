package models

import "time"

type Review struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	AuthorID  uint      `json:"author_id"`
	ProductID string      `json:"product_id"`
	Rate      int8      `json:"rate"`
	Comment   string    `json:"comment"`
}

type AvgAndCount struct {
	AvgRate      float64 `json:"average_rate"`
	NumberReview int     `json:"count_reviews"`
}
