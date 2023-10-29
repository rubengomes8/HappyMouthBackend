package users

type User struct {
	ID       int64  `json:"id,omitempty" gorm:"primaryKey:type:int"`
	Username string `json:"username,omitempty" gorm:"username"`
	Passhash string `json:"passhash,omitempty" gorm:"passhash"`
	Email    string `json:"email,omitempty" gorm:"email"`
}

type UserFilters struct {
	ID       int64  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}
