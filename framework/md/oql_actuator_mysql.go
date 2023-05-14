package md

import (
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/utils"
)

func init() {
	RegisterOQLActuator(&mysqlActuator{})
}

// 公共查询
type mysqlActuator struct {
	from   string
	offset int
	limit  int
}

func (mysqlActuator) GetName() string {
	return utils.ORM_DRIVER_MYSQL
}
func (s *mysqlActuator) Count(oql OQL, out interface{}) OQL {
	oql.Parse()
	q:=db.Default().DB
	if statement := oql.BuildFrom(); statement.Affected > 0 {
		q=q.Table(statement.Query)
	}
	if statement := oql.BuildWheres(); statement.Affected > 0 {
		q=q.Where(statement.Query,statement.Args...)
	}
	if statement := oql.BuildGroups(); statement.Affected > 0 {
		q=q.Group(statement.Query)
	}
	if statement := oql.BuildHaving(); statement.Affected > 0 {
		q=q.Having(statement.Query,statement.Args...)
	}
	if err := q.Count(out).Error; err != nil {
		oql.AddErr(err)
	}
	return oql
}
func (s *mysqlActuator) Pluck(oql OQL, column string, value interface{}) OQL {
	return oql
}
func (s *mysqlActuator) Take(oql OQL, out interface{}) OQL {
	if err := db.Default().Take(out).Error; err != nil {
		oql.AddErr(err)
	}
	return oql
}
func (s *mysqlActuator) Find(oql OQL, out interface{}) OQL {
	oql.Parse()
	q:=db.Default().DB
	if statement := oql.BuildFrom(); statement.Affected > 0 {
		q=q.Table(statement.Query)
	}
	if statement := oql.BuildSelects(); statement.Affected > 0 {
		q=q.Select(statement.Query,statement.Args...)
	}
	if statement := oql.BuildWheres(); statement.Affected > 0 {
		q=q.Where(statement.Query,statement.Args...)
	}
	if statement := oql.BuildGroups(); statement.Affected > 0 {
		q=q.Group(statement.Query)
	}
	if statement := oql.BuildHaving(); statement.Affected > 0 {
		q=q.Having(statement.Query,statement.Args...)
	}
	if err := q.Find(out).Error; err != nil {
		oql.AddErr(err)
	}
	return oql
}
func (s *mysqlActuator) Create(oql OQL, data interface{}) OQL {
	return oql
}
func (s *mysqlActuator) Update(oql OQL, data interface{}) OQL {
	return oql
}
func (s *mysqlActuator) Delete(oql OQL) OQL {
	return oql
}
