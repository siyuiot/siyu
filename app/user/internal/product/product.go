package product

import (
	"time"
)

type Product struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Des        string    `json:"des"`
	Status     int       `json:"status"`
	CategoryId int       `json:"category_id"`
	BeginAt    time.Time `json:"begin_at"`
	ExpiredAt  time.Time `json:"expired_at"`
	Img        string    `json:"img"`
	Creator    int       `json:"creator"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Product) TableName() string {
	return "product"
}

func GetProductById(id int) (res *Product, err error) {
	err = this.DBORM.Model(Product{}).Where("id = ?", id).Find(res).Error
	return
}
