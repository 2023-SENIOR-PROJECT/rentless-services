// >> Unfinished code << 

package models

import "time"

type User struct {
	ID        uint      `db:"id" gorm:"primaryKey"`
	Firstname string    `db:"firstname"`
	Lastname  string    `db:"lastname"`
	Age       int       `db:"age"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at" gorm:"not null"`
}
