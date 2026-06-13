package model

type User struct {
	ID    uint64 `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
	Age   uint   `db:"age" json:"age"`
}
