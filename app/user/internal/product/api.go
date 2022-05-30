package product

// 产品：用户服务

import (
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/siyuiot/siyu/pkg/qlog"
)

var this *object

var orderDB *gorm.DB

type Object interface {
	QueryProductSkuByPid(pid int) ([]ProductSku, error)
	QueryProductSkuBySkuId(skuId int) *ProductSku
}

type Option struct {
	Log   *qlog.Entry
	Db    *sql.DB
	DbRo  *sql.DB
	DBORM *gorm.DB
}

type object struct {
	Option
}

func New(o Option) {
	this = &object{
		Option: o,
	}
}

func Instance() *object {
	return this
}

func (o *object) QueryProductSkuByPid(id int) (res []ProductSku, err error) {
	if id <= 0 {
		err = fmt.Errorf("param err")
		return
	}
	err = o.DBORM.Table(ProductSku{}.TableName()+" AS ps").
		Select("ps.*").
		Joins("LEFT JOIN product AS p ON ps.product_id = p.id").
		Where("p.id = ?", id).
		Order("ps.price asc").
		Scan(&res).Error
	return
}

func (o *object) QueryProductSkuBySkuId(skuId int) *ProductSku {
	r := new(ProductSku)
	var qstr string
	if skuId > 0 {
		qstr += fmt.Sprintf(" and ps.id = %d ", skuId)
	}
	sqlstr := `
	select
	coalesce(ps.id,0) as id,
	coalesce(ps.name,'') as name,
	coalesce(ps.price,0) as price,
	coalesce(ps.price_origin,0) as price_origin,
	coalesce(ps.status,0) as status,
	coalesce(ps.img,'') as img,
	coalesce(ps.des,'') as des,
	coalesce(ps.created_at,'0001-01-01') as created_at,
	coalesce(ps.updated_at,'0001-01-01') as updated_at
	from product_sku ps
	left join product p on ps.product_id = p.id
	where 1=1
	`
	sqlstr += qstr
	// app.Log.Debug(sqlstr)
	err := o.DbRo.QueryRow(sqlstr).Scan(&r.Id, &r.Name, &r.Price, &r.PriceOrigin, &r.Status, &r.Img, &r.Des, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		o.Log.Errorf("param=%d,sql=%s,err=%v", skuId, sqlstr, err)
		return nil
	}
	return r
}
