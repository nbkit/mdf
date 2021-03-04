package rules

import (
	"github.com/nbkit/mdf/framework/files"
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
	"sort"
)

type entityImportPre struct {
}

func newEntityImportPre() *entityImportPre {
	return &entityImportPre{}
}
func (s *entityImportPre) Register() md.RuleRegister {
	return md.RuleRegister{Code: "importPre", OwnerType: md.RuleType_Widget, OwnerCode: "md"}
}
func (s *entityImportPre) Exec(flow *utils.FlowContext) {
	if flow.Request.Data == nil {
		flow.Error("没有要导入的数据")
		return
	}
	if items, ok := flow.Request.Data.([]files.ImportData); !ok {
		flow.Error("导入的数据非法！")
		return
	} else {
		if err := s.batchImport(items); err != nil {
			flow.Error(err)
			return
		}
	}
}

func (s *entityImportPre) batchImport(datas []files.ImportData) error {
	if len(datas) > 0 {
		nameList := make(map[string]int)
		nameList["Entity"] = 1
		nameList["Props"] = 2

		sort.Slice(datas, func(i, j int) bool { return nameList[datas[i].SheetName] < nameList[datas[j].SheetName] })

		entities := make([]md.MDEntity, 0)
		fields := make([]md.MDField, 0)
		for _, item := range datas {
			if item.SheetName == "Entity" {
				if d, err := s.toEntities(item); err != nil {
					return err
				} else if len(d) > 0 {
					entities = append(entities, d...)
				}
			}
			if item.SheetName == "Props" {
				if d, err := s.toFields(item); err != nil {
					return err
				} else if len(d) > 0 {
					fields = append(fields, d...)
				}
			}
		}
		if len(entities) > 0 {
			for i, entity := range entities {
				for _, field := range fields {
					if entity.ID == field.EntityID {
						if entities[i].Fields == nil {
							entities[i].Fields = make([]md.MDField, 0)
						}
						entities[i].Fields = append(entities[i].Fields, field)
					}
				}
			}
			md.MDSv().AddEntities(entities)
		}
	}
	return nil
}
func (s *entityImportPre) toEntities(data files.ImportData) ([]md.MDEntity, error) {
	if len(data.Data) == 0 {
		return nil, nil
	}
	items := make([]md.MDEntity, 0)
	for _, row := range data.Data {
		item := md.MDEntity{}
		item.ID = files.GetCellValue("ID", row)
		item.Name = files.GetCellValue("Name", row)
		item.Type = files.GetCellValue("Type", row)
		item.TableName = files.GetCellValue("TableName", row)
		item.Domain = files.GetCellValue("Domain", row)
		item.System = utils.ToSBool(files.GetCellValue("System", row))
		items = append(items, item)
	}
	return items, nil
}
func (s *entityImportPre) toFields(data files.ImportData) ([]md.MDField, error) {
	if len(data.Data) == 0 {
		return nil, nil
	}
	items := make([]md.MDField, 0)
	for _, row := range data.Data {
		item := md.MDField{}
		if cValue := files.GetCellValue("EntityID", row); cValue != "" {
			item.EntityID = cValue
		}
		if cValue := files.GetCellValue("Name", row); cValue != "" {
			item.Name = cValue
		}
		if cValue := files.GetCellValue("Code", row); cValue != "" {
			item.Code = cValue
		}
		if cValue := files.GetCellValue("TypeID", row); cValue != "" {
			item.TypeID = cValue
		}
		if cValue := files.GetCellValue("Kind", row); cValue != "" {
			item.Kind = cValue
		}
		if cValue := files.GetCellValue("ForeignKey", row); cValue != "" {
			item.ForeignKey = cValue
		}
		if cValue := files.GetCellValue("AssociationKey", row); cValue != "" {
			item.AssociationKey = cValue
		}
		if cValue := files.GetCellValue("DbName", row); cValue != "" {
			item.DbName = cValue
		}
		if cValue := files.GetCellValue("DbName", row); cValue != "" {
			item.DbName = cValue
		}
		if cValue := utils.ToInt(files.GetCellValue("Length", row)); cValue >= 0 {
			item.Length = cValue
		}
		if cValue := utils.ToInt(files.GetCellValue("Precision", row)); cValue >= 0 {
			item.Precision = cValue
		}
		if cValue := files.GetCellValue("DefaultValue", row); cValue != "" {
			item.DefaultValue = cValue
		}
		if cValue := files.GetCellValue("MaxValue", row); cValue != "" {
			item.MaxValue = cValue
		}
		if cValue := files.GetCellValue("MinValue", row); cValue != "" {
			item.MinValue = cValue
		}
		if cValue := files.GetCellValue("Tags", row); cValue != "" {
			item.Tags = cValue
		}
		if cValue := files.GetCellValue("Limit", row); cValue != "" {
			item.Limit = cValue
		}
		item.Nullable = utils.ToSBool(files.GetCellValue("Nullable", row))
		item.IsPrimaryKey = utils.ToSBool(files.GetCellValue("IsPrimaryKey", row))
		items = append(items, item)
	}
	return items, nil
}
