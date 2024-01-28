package coins

import "time"

type UserCoins struct {
	UserID    int64      `json:"user_id,omitempty" gorm:"user_id,primaryKey:type:int"`
	Coins     int        `json:"coins,omitempty" gorm:"coins"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"updated_at"`
}
