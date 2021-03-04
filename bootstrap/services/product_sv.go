package services

import (
	"github.com/nbkit/mdf/db"
	"sync"

	"github.com/nbkit/mdf/bootstrap/model"
)

type IProductSv interface {
}

type productSvImpl struct {
	*sync.Mutex
}

func ProductSv() IProductSv {
	return productSv
}

var productSv IProductSv = newProductSv()

/**
* 创建服务实例
 */
func newProductSv() *productSvImpl {
	return &productSvImpl{Mutex: &sync.Mutex{}}
}

func (s *productSvImpl) GetServiceByCode(entID, idOrCode string) (*model.ProductService, error) {
	old := model.ProductService{}
	if err := db.Default().Where("ent_id=? and (code=? or id=?)", entID, idOrCode, idOrCode).Take(&old).Error; err != nil {
		return nil, err
	}
	return &old, nil
}
