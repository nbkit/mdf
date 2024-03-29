package utils

import (
	"reflect"
)

func To(v, to interface{}, tags string) (interface{}, error) {
	if nil == v || nil == to {
		return nil, errInvalidParams
	}
	return ToT(v, reflect.ValueOf(to).Type(), tags)
}

func ToT(v interface{}, t reflect.Type, tags string) (interface{}, error) {
	if nil == v || nil == t {
		return nil, errInvalidParams
	}
	vl, err := ToType(reflect.ValueOf(v), t, tags)
	return vl, err
}

func ToType(v reflect.Value, t reflect.Type, tags string) (interface{}, error) {
	if reflect.Ptr == v.Kind() {
		v = ReflectTarget(v)
	}

	var err error
	switch t.Kind() {
	case reflect.String:
		return ToString(v.Interface()), nil
	case reflect.Bool:
		return ToBool(v.Interface()), nil
	case reflect.Int:
		return ToInt(v.Interface()), nil
	case reflect.Int8:
		return (int8)(ToInt(v.Interface())), nil
	case reflect.Int16:
		return (int16)(ToInt(v.Interface())), nil
	case reflect.Int32:
		return ToInt32(v.Interface()), nil
	case reflect.Int64:
		return ToInt64(v.Interface()), nil
	case reflect.Uint:
		return ToUint(v.Interface()), nil
	case reflect.Uint8:
		return (uint8)(ToUint(v.Interface())), nil
	case reflect.Uint16:
		return (uint16)(ToUint(v.Interface())), nil
	case reflect.Uint32:
		return ToUint32(v.Interface()), nil
	case reflect.Uint64:
		return ToUint64(v.Interface()), nil
	case reflect.Float32:
		return ToFloat32(v.Interface()), nil
	case reflect.Float64:
		return ToFloat64(v.Interface()), nil
	case reflect.Slice:
		slice := reflect.New(t)
		if err = ToSlice(slice.Interface(), v.Interface(), tags); nil == err {
			return slice.Elem().Interface(), nil
		}
		break
	case reflect.Map:
		mp := reflect.MakeMap(t)
		if err = ToMap(mp.Interface(), v.Interface(), tags, true); nil == err {
			return mp.Interface(), nil
		}
		break
	case reflect.Ptr:
		var vl interface{}
		if vl, err = ToType(v, t.Elem(), tags); nil == err {
			if reflect.Ptr != t.Elem().Kind() {
				return vl, nil
			} else {
				vlPtr := reflect.Zero(t)
				val := reflect.ValueOf(vl)
				vlPtr.Set(val)
				return vlPtr.Interface(), nil
			}
		}
		break
	case reflect.Struct:
		st := reflect.New(t).Interface()
		if err = ToStruct(st, v.Interface(), tags); nil == err {
			return st, nil
		}
		break
	}
	return nil, err
}
