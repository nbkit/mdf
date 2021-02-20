package services

import (
	"github.com/nbkit/mdf/bootstrap/model"
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/utils"
)

type ITokenSv interface {
}

type tokenSvImpl struct{}

var tokenSv ITokenSv = newTokenSvImpl()

func TokenSv() ITokenSv {
	return tokenSv
}

/**
* 创建服务实例
 */
func newTokenSvImpl() *tokenSvImpl {
	return &tokenSvImpl{}
}

func (s *tokenSvImpl) Create(token *model.AuthToken) *model.AuthToken {
	token.ID = utils.GUID()
	token.Token = utils.GUID()
	db.Default().Create(token)
	return token
}
func (s *tokenSvImpl) Get(token string) (*model.AuthToken, bool) {
	t := &model.AuthToken{}
	if err := db.Default().Model(&model.AuthToken{}).Where("token=?", token).Take(t).Error; err != nil {
		return t, false
	}
	return t, true
}

func (s *tokenSvImpl) Delete(ids []string) error {
	if err := db.Default().Delete(model.AuthToken{}, "id in (?)", ids).Error; err != nil {
		return err
	}
	return nil
}
func (s *tokenSvImpl) GetAndUse(token string) (*model.AuthToken, bool) {
	t := &model.AuthToken{}
	if err := db.Default().Model(&model.AuthToken{}).Where("token=?", token).Take(t).Error; err != nil {
		return t, false
	}
	db.Default().Where("id = ?", t.ID).Delete(&model.AuthToken{})
	return t, true
}
