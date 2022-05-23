package product

import (
	"time"
)

type MallAttribute struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (MallAttribute) TableName() string {
	return "mall_attribute"
}
