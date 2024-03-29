package utils

import (
	"database/sql"
	"reflect"
	"strings"
	"time"
)

func ToStruct(dst, src interface{}, tag string) (err error) {
	if nil == dst || nil == src {
		err = errInvalidParams
		return
	}

	if sintf, ok := dst.(sql.Scanner); ok {
		if nil == sintf.Scan(src) {
			return
		}
	}

	switch dst.(type) {
	case time.Time, *time.Time:
		switch v := src.(type) {
		case time.Time:
			s := ReflectTarget(reflect.ValueOf(dst))
			s.Set(reflect.ValueOf(v))
			break
		case *time.Time:
			s := ReflectTarget(reflect.ValueOf(dst))
			s.Set(reflect.ValueOf(*v))
			break
		case string:
			var tm time.Time
			if tm, err = ParseTime(v); nil == err {
				s := ReflectTarget(reflect.ValueOf(dst))
				s.Set(reflect.ValueOf(tm))
			}
			break
		case int64:
			s := ReflectTarget(reflect.ValueOf(dst))
			s.Set(reflect.ValueOf(time.Unix(v, 0)))
			break
		default:
			err = errUnsupportedType
		}
		break
	default:

		s := ReflectTarget(reflect.ValueOf(dst))
		t := s.Type()

		switch src.(type) {
		case map[interface{}]interface{}, map[string]interface{}, map[string]string:
			for i := 0; i < s.NumField(); i++ {
				f := s.Field(i)
				if f.CanSet() {
					// Get passable field names
					names := fieldNames(t.Field(i), tag)
					if nil == names || len(names) < 1 {
						continue
					}

					// Get value from map
					v := mapValueByStringKeys(src, names)

					// Set field value
					if nil == v {
						f.Set(reflect.Zero(f.Type()))
					} else {
						switch f.Kind() {
						case reflect.Struct:
							if err = ToStruct(f.Addr().Interface(), v, tag); nil != err {
								return
							}
							break
						default:
							var vl interface{}
							if vl, err = ToType(reflect.ValueOf(v), f.Type(), tag); nil == err {
								val := reflect.ValueOf(vl)
								if val.Kind() == reflect.Ptr && val.Kind() != f.Kind() {
									val = val.Elem()
								}
								if val.Kind() == f.Kind() {
									f.Set(val)
								} else {
									err = errUnsupportedType
									break
								}
							} else {
								return
							}
						} // end switch
					} // end else
				}
			}
			break
		default:
			err = errUnsupportedType
		}
	}

	return
}

func StructFields(st interface{}, tag string) []string {
	fields := []string{}

	s := ReflectTarget(reflect.ValueOf(st))
	t := s.Type()

	for i := 0; i < s.NumField(); i++ {
		fname, _ := fieldName(t.Field(i), tag)
		if len(fname) > 0 && "-" != fname {
			fields = append(fields, fname)
		}
	}
	return fields
}

func StructFieldTags(st interface{}, tag string) map[string]string {
	fields := map[string]string{}
	keys, values := StructFieldTagsUnsorted(st, tag)

	for i, k := range keys {
		fields[k] = values[i]
	}
	return fields
}

func StructFieldTagsUnsorted(st interface{}, tag string) ([]string, []string) {
	keys := []string{}
	values := []string{}

	s := ReflectTarget(reflect.ValueOf(st))
	t := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := t.Field(i)
		tag := fieldTag(f, tag)
		if len(tag) > 0 && "-" != tag {
			keys = append(keys, f.Name)
			values = append(values, tag)
		}
	}
	return keys, values
}

///////////////////////////////////////////////////////////////////////////////
/// MARK: Helpers
///////////////////////////////////////////////////////////////////////////////

var fieldNameArr = []string{"field", "schema", "sql", "json", "xml", "yaml"}

func fieldNames(f reflect.StructField, tag string) []string {
	names := fieldTagArr(f, tag)
	switch names[0] {
	case "", "-":
		return []string{f.Name}
		break
	default:
	}
	return []string{names[0], f.Name}
}

func fieldName(f reflect.StructField, tag string) (name string, omitempty bool) {
	names := fieldTagArr(f, tag)
	name = names[0]
	if len(names) > 1 {
		if "omitempty" == names[len(names)-1] {
			omitempty = true
		}
	}
	if "" == name {
		name = f.Name
	}
	return
}

func fieldTagArr(f reflect.StructField, tag string) []string {
	return strings.Split(fieldTag(f, tag), ",")
}

func fieldTag(f reflect.StructField, tag string) string {
	if "-" != tag {
		var fields string
		var tags []string

		if len(tag) > 0 {
			tags = strings.Split(tag, ",")
		} else {
			tags = fieldNameArr
		}

		for _, k := range tags {
			fields = f.Tag.Get(k)
			if len(fields) > 0 {
				break
			}
		}

		if len(fields) > 0 {
			if "-" == fields {
				return ""
			}
			return fields
		}
	}
	return f.Name
}
