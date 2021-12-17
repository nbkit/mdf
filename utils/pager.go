package utils

import "reflect"

type Pager struct {
	Value    interface{}
	Error    error
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	LastPage int `json:"last_page"`
	Items    int `json:"items"` //当前条数
	Total    int `json:"total"` //总记录数
}

func NewPager(value interface{}, page int, pageSize int, itemTotal int) *Pager {
	item := Pager{Value: value, Page: page, PageSize: pageSize, Total: itemTotal}
	if item.Total > 0 && item.PageSize > 0 {
		item.LastPage = item.Total / item.PageSize
	}
	if value != nil {
		if aa := reflect.TypeOf(value); aa != nil {
			item.Items = reflect.ValueOf(value).Len()
		}
	}
	return &item
}
func (c *Pager) ToResource() Map {
	rtn := Map{"data": nil}
	if c == nil {
		return rtn
	}
	rtn["pager"] = ToPagerItem(c)
	rtn["data"] = c.Value
	return rtn
}

type PagerRes struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	LastPage int `json:"last_page"`
	Items    int `json:"items"` //当前条数
	Total    int `json:"total"` //总记录数
}

func (c *PagerRes) ToResource(obj *Pager) Map {
	rtn := Map{"data": nil}
	if obj == nil {
		return rtn
	}
	rtn["pager"] = ToPagerItem(obj)
	rtn["data"] = obj.Value
	return rtn
}

func ToPagerRes(obj *Pager) Map {
	return ToPagerItem(obj).ToResource(obj)
}
func ToPagerItem(obj *Pager) *PagerRes {
	return &PagerRes{
		Page:     obj.Page,
		PageSize: obj.PageSize,
		Total:    obj.Total,
		Items:    obj.Items,
		LastPage: obj.LastPage,
	}
}
