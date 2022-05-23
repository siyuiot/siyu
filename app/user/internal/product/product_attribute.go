package product

import (
	"time"
)

type ProductAttribute struct {
	Id        int       `json:"id"`
	ProductId int       `json:"product_id"`
	AttrId    int       `json:"attr_id"`
	Seq       int       `json:"seq"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ProductAttribute) TableName() string {
	return "product_attribute"
}
