package coins

import "time"

type UserCoins struct {
	UserID    int64      `json:"user_id,omitempty" gorm:"user_id,primaryKey:type:int"`
	Coins     int        `json:"coins,omitempty" gorm:"coins"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"updated_at"`
}

func (uc UserCoins) Subtract(qty int) (UserCoins, error) {
	if qty > uc.Coins {
		return UserCoins{}, ErrNotEnoughCoins
	}
	now := time.Now().UTC()
	return UserCoins{
		UserID:    uc.UserID,
		Coins:     uc.Coins - qty,
		CreatedAt: uc.CreatedAt,
		UpdatedAt: &now,
	}, nil
}
