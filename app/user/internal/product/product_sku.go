package product

import (
	"fmt"
	"time"
)

type ProductSku struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Price       int       `json:"price"`
	PriceOrigin int       `json:"price_origin"`
	ProductId   int       `json:"product_id"`
	Status      int       `json:"status"`
	Img         string    `json:"img"`
	Des         string    `json:"des"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ProductSku) TableName() string {
	return "product_sku"
}

type SkuValues struct {
	SkuId    int    `json:"-"`
	AttrName string `json:"attr_name"`
	AttrId   string `json:"attr_id"`
	ValId    string `json:"val_id"`
	ValName  string `json:"val_name"`
}

type SkuInfo struct {
	ProductSku
	Values []SkuValues `json:"values"`
}

func GetProductSkuInfoById(id int) (res SkuInfo, err error) {
	data, err := GetProductSkuById(id)
	if err != nil {
		return
	}
	res.ProductSku = data
	values, err := GetProductSkuValues(id)
	if err != nil {
		return
	}
	res.Values = values
	return
}

func GetProductSkuByPid(id int) (res []ProductSku, err error) {
	if id <= 0 {
		err = fmt.Errorf("param err")
		return
	}
	err = this.DBORM.Table(ProductSku{}.TableName()+" AS ps").
		Select("ps.*").
		Joins("LEFT JOIN product AS p ON ps.product_id = p.id").
		Where("p.id = ?", id).
		Order("ps.price asc").
		Scan(&res).Error
	return
}

func GetProductSkuById(id int) (res ProductSku, err error) {
	if id <= 0 {
		err = fmt.Errorf("param err")
		return
	}
	err = this.DBORM.Model(ProductSku{}).Where("id = ?", id).Find(&res).Error
	return
}

func GetProductSkuValues(skuId int) (res []SkuValues, err error) {
	err = this.DBORM.Table(ProductSku{}.TableName()+" AS ps").
		Select("ps.id as sku_id,ma.name AS attr_name,ma.id AS attr_id,mv.id AS val_id,mv.name AS val_name").
		Joins("JOIN product_sku_val AS psv ON psv.sku_id = ps.id ").
		Joins("JOIN mall_value AS mv ON mv.id = psv.val_id").
		Joins("JOIN mall_attribute AS ma ON ma.id = mv.attr_id ").
		Where("ps.id = ?", skuId).
		Order("ma.seq ASC,ma.id ASC").
		Scan(&res).Error
	return
}

// int slice 去重
func UniqueIntSlice(input []int) (u []int) {
	if len(input) == 0 {
		return
	}
	u = make([]int, 0, len(input))
	m := make(map[int]struct{})
	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = struct{}{}
			u = append(u, val)
		}
	}
	return
}

func GetProductSkuValuesBySkuList(skuIdList []int) (res map[int][]SkuValues, err error) {
	skuIdList = UniqueIntSlice(skuIdList)
	list := make([]SkuValues, 0)
	err = this.DBORM.Table(ProductSku{}.TableName()+" AS ps").
		Select("ps.id as sku_id,ma.name AS attr_name,ma.id AS attr_id,mv.id AS val_id,mv.name AS val_name").
		Joins("JOIN product_sku_val AS psv ON psv.sku_id = ps.id ").
		Joins("JOIN mall_value AS mv ON mv.id = psv.val_id").
		Joins("JOIN mall_attribute AS ma ON ma.id = mv.attr_id ").
		Where("ps.id in (?)", skuIdList).
		Order("ma.seq ASC,ma.id ASC").
		Scan(&list).Error
	res = make(map[int][]SkuValues, 0)
	for _, values := range list {
		res[values.SkuId] = append(res[values.SkuId], values)
	}
	return
}
