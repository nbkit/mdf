package services

import (
	"sync"
)

type IRoleSv interface {
}
type roleSvImpl struct {
	*sync.Mutex
}

var roleSv IRoleSv = newRoleSvImpl()

func RoleSv() IRoleSv {
	return roleSv
}

/**
* 创建服务实例
 */
func newRoleSvImpl() *roleSvImpl {
	return &roleSvImpl{Mutex: &sync.Mutex{}}
}
