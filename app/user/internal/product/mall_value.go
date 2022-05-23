package product

import (
	"time"
)

type MallValue struct {
	Id        int             `json:"id"`
	AttrId    int             `json:"attr_id"`
	Name      string          `json:"name"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	MallAttrs []MallAttribute `json:"mall_attrs" gorm:"foreignkey:id;association_foreignkey:attr_id"`
}

func (MallValue) TableName() string {
	return "mall_value"
}
