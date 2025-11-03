package model

// DBUser – ...
type DBUser struct {
	ID           int64  `json:"id" db:"id"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"pass_hash" db:"pass_hash"`
}

// User – ...
type User struct {
	ID    int64
	Email string
}

// App – ...
type App struct {
	ID     int64
	Name   string
	Secret string
}
