package product

import (
	"time"
)

type ProductSkuVal struct {
	Id        int       `json:"id"`
	SkuId     int       `json:"sku_id"`
	ValId     int       `json:"val_id"`
	Seq       int       `json:"seq"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联模型
	MallValues []MallValue `json:"mall_values" gorm:"foreignkey:id;association_foreignkey:val_id"`
}

func (ProductSkuVal) TableName() string {
	return "product_sku_val"
}

func GetProductSkuValBySkuId(id int) (list []ProductSkuVal, err error) {
	err = this.DBORM.Model(ProductSkuVal{}).
		Preload("MallValues").
		Preload("MallValues.MallAttrs").
		Where("sku_id = ?", id).Find(&list).Error
	return
}
