package utils

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/nbkit/mdf/log"
)

// JSONTime format json time field by myself
type SJson struct {
	value      interface{}
	jsonString string
	valid      bool // Valid is true if jsonString is not NULL
	isMap      bool
}

var SJson_Null = SJson{value: nil, jsonString: "", valid: false}

func ToSJson(val interface{}) SJson {
	s := SJson{}
	s.Parse(val)
	return s
}
func (t *SJson) Parse(val interface{}) SJson {
	if v, err := json.Marshal(val); err != nil {
		t.valid = false
		t.jsonString = ""
		t.value = nil
	} else {
		t.value = val
		t.jsonString = string(v)
		t.valid = true
	}
	t.formatMapValue()
	return *t
}

func (t *SJson) formatMapValue() {
	if t.valid && len(t.jsonString) >= 2 && t.jsonString[0] == '{' && t.jsonString[len(t.jsonString)-1] == '}' {
		mapValue := make(map[string]interface{})
		if err := json.Unmarshal([]byte(t.jsonString), &mapValue); err != nil {
			t.isMap = false
		}
		t.value = mapValue
		t.isMap = true
	} else {
		t.isMap = false
	}

}
func (t SJson) String() string {
	return t.jsonString
}
func (t SJson) Valid() bool {
	return t.valid
}
func (t SJson) Equal(v SJson) bool {
	return t.jsonString == v.jsonString
}

func (t SJson) IsArray() bool {
	if t.jsonString == "" {
		return false
	}
	return t.jsonString[0] == '['
}

func (t SJson) NotEqual(v SJson) bool {
	return t.jsonString != v.jsonString
}
func (t SJson) MarshalJSON() ([]byte, error) {
	if !t.valid {
		return []byte("null"), nil
	}
	return []byte(t.String()), nil
}
func (t *SJson) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		t.value = nil
		t.jsonString = ""
		t.valid = false
		return nil
	}
	strData := string(data)
	if strData == "" || strData == "null" || strData == " " || strData == "undefined" {
		t.value = nil
		t.jsonString = ""
		t.valid = false
		return nil
	}
	var vv interface{}
	if err := json.Unmarshal(data, &vv); err != nil {
		t.value = nil
		t.jsonString = ""
		t.valid = false
	} else {
		t.value = vv
		t.jsonString = strData
		t.valid = t.jsonString != ""
	}
	t.formatMapValue()
	return nil
}

// Value implements the driver Valuer interface.
func (t SJson) Value() (driver.Value, error) {
	if !t.valid {
		return nil, nil
	}
	return t.jsonString, nil
}

// Scan implements the Scanner interface.
func (t *SJson) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	jsonStr := ToString(v)
	if jsonStr == "" {
		t.value = nil
		t.jsonString = jsonStr
		t.valid = t.jsonString != ""
		return nil
	}
	var vv interface{}
	if err := json.Unmarshal([]byte(jsonStr), &vv); err != nil {
		if err != nil {
			log.Print("unmarshal Json string failed", err)
		}
		if strValue, e := json.Marshal(jsonStr); e == nil {
			t.value = jsonStr
			t.jsonString = string(strValue)
			t.valid = t.jsonString != ""
		}
	} else {
		t.value = vv
		t.jsonString = jsonStr
		t.valid = t.jsonString != ""
	}
	t.formatMapValue()
	return nil
}
func (t SJson) OrmDataType(driver string) string {
	return "text"
}
func (t SJson) GetObject(obj interface{}) error {
	return json.Unmarshal([]byte(t.jsonString), &obj)
}
func (t SJson) GetString() string {
	return ToString(t.value)
}
func (t SJson) GetValue() interface{} {
	return t.value
}
func (t SJson) GetInterfaceSlice() []interface{} {
	return ToInterfaceSlice(t.value)
}
func (t SJson) GetMapValueBy(key string) (interface{}) {
	if t.isMap && t.value != nil {
		if mapValue, isMapValue := t.value.(map[string]interface{}); isMapValue {
			if mv, ok := mapValue[key]; ok {
				return mv
			}
		}
	}
	return nil
}
func (t SJson) GetMapStringValueBy(key string) string {
	if t.isMap && t.value != nil {
		if mapValue, isMapValue := t.value.(map[string]interface{}); isMapValue {
			if mv, ok := mapValue[key]; ok {
				return ToString(mv)
			}
		}
	}
	return ""
}
func (t *SJson) SetMapValueBy(key string, value interface{}) bool {
	if t.isMap {
		if t.value == nil {
			t.value = make(map[string]interface{})
		}
		if mapValue, isMapValue := t.value.(map[string]interface{}); isMapValue {
			mapValue[key] = value
			t.value = mapValue
			t.Parse(t.value)
			return true
		}
	}
	return false
}
