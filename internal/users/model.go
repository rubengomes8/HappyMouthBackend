package users

import "time"

type User struct {
	ID        int64      `json:"id,omitempty" gorm:"primaryKey:type:int"`
	Username  string     `json:"username,omitempty" gorm:"username"`
	Passhash  string     `json:"passhash,omitempty" gorm:"passhash"`
	Email     string     `json:"email,omitempty" gorm:"email"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"deleted_at"`
}

type UserFilters struct {
	ID       int64  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
