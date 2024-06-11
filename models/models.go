package models

type Employee struct {
	ID       int     `json:"id" db:"id"`
	Name     string  `json:"name" db:"name"`
	Position string  `json:"position" db:"position"`
	Salary   float64 `json:"salary" db:"salary"`
}
