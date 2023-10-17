package models

import "time"

type Review struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	AuthorID  uint      `json:"author_id"`
	ProductID uint      `json:"product_id"`
	Rate      int8      `json:"rate"`
	Comment   string    `json:"comment"`
}
