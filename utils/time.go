package utils

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/nbkit/mdf/log"
)

// JSONTime format json time field by myself
type Time struct {
	time.Time
}

const (
	Layout_YYYYMM          = "2006-01"
	Layout_YYYYMMDD        = "2006-01-02"
	Layout_YYYYMMDDHHIISS  = "2006-01-02 15:04:05"
	Layout_YYYYMMDDHHIISST = "2006-01-02T15:04:05"
	Layout_YYYYMMDD2       = "20060102"
)

// yyyy-MM-dd HH:mm:ss.SSSSS
// 2006-01-02T15:04:05.000000
func formatStrEnd(from int, in string, target rune) (string, int) {
	var out = new(strings.Builder)
	for i := from; i < len(in); i++ {
		r := rune(in[i])
		from = i
		if r == target {
			out.WriteRune(r)
			if i == len(in)-1 {
				return out.String(), i + 1
			}
			continue
		}
		return out.String(), i
	}
	return "", from + 1
}

// y:年,M:年中的月份
// w:年中的周数	,W:月份中的周数
// D:年中的天数,d:月份中的天数
// H:一天中的小时数（0-23）,m:小时中的分钟数	,s:分钟中的秒数,S:毫秒数
// yyyy-MM-dd HH:mm:ss.SSSSS
// 2006-01-02T15:04:05.000000
func TimeFormatStr(layout string) string {
	var i = 0
	var t = new(strings.Builder)
	for i < len(layout) {
		c := layout[i]
		switch c {
		case 'y': // 年[year]
			y, endIndex := formatStrEnd(i, layout, 'y')
			if length := len(y); length > 3 {
				t.WriteString("2006")
			} else {
				t.WriteString("06")
			}
			i = endIndex
		case 'M': // 月[month]
			_, endIndex := formatStrEnd(i, layout, 'M')
			t.WriteString("01")
			i = endIndex
		case 'd': // 月份中的天数[number]
			_, endIndex := formatStrEnd(i, layout, 'd')
			t.WriteString("02")
			i = endIndex
		case 'H': // 一天中的小时数，0-23[number]
			_, endIndex := formatStrEnd(i, layout, 'H')
			t.WriteString("15")
			i = endIndex
		case 'm': // 小时中的分钟数[number]
			_, endIndex := formatStrEnd(i, layout, 'm')
			t.WriteString("04")
			i = endIndex
		case 's': // 分钟中的秒数[number]
			_, endIndex := formatStrEnd(i, layout, 's')
			t.WriteString("05")
			i = endIndex
		case 'S': // 毫秒数[number]
			ss, endIndex := formatStrEnd(i, layout, 'S')
			t.WriteString(fmt.Sprintf("%0*d", len(ss), 0))
			i = endIndex
		default:
			t.WriteByte(c)
			i = i + 1
		}
	}
	return t.String()
}
func TimeNow() Time {
	return ToTime(time.Now())
}
func TimeNowPtr() *Time {
	return ToTimePtr(time.Now())
}
func ToTime(value interface{}) Time {
	if value == "" || value == nil {
		return Time{}
	}
	if v, ok := value.(time.Time); ok {
		return Time{v}
	}
	if v, ok := value.(string); ok {
		v = strings.Replace(v, `"`, "", -1)
		layout := Layout_YYYYMMDDHHIISS
		data := []rune(v)
		if len(data) >= len(Layout_YYYYMMDDHHIISST) && strings.Contains(v, "T") {
			layout = Layout_YYYYMMDDHHIISST
			data = data[:len(Layout_YYYYMMDDHHIISST)]
		} else if len(data) >= len(Layout_YYYYMMDDHHIISS) {
			layout = Layout_YYYYMMDDHHIISS
			data = data[:len(Layout_YYYYMMDDHHIISS)]
		} else if len(data) >= len(Layout_YYYYMMDD) {
			layout = Layout_YYYYMMDD
			data = data[:len(Layout_YYYYMMDD)]
		} else if len(data) >= len(Layout_YYYYMM) {
			layout = Layout_YYYYMM
			data = data[:len(Layout_YYYYMM)]
		}
		now, err := time.ParseInLocation(layout, string(data), time.Local)
		if err != nil {
			log.Print("string to ToTime failed", err)
		}
		return Time{now}
	}
	return Time{time.Now()}
}
func ToTimePtr(value interface{}) *Time {
	t := ToTime(value)
	return &t
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t Time) MarshalJSON() ([]byte, error) {
	if t.UnixNano() <= 0 || t.Unix() <= 0 || t.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}
func (t *Time) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		*t = Time{}
		return nil
	}
	ds := strings.Replace(string(data), `"`, "", -1)
	layout := Layout_YYYYMMDDHHIISS
	data = []byte(ds)
	if len(data) >= len(Layout_YYYYMMDDHHIISST) && strings.Contains(ds, "T") {
		layout = Layout_YYYYMMDDHHIISST
		data = data[:len(Layout_YYYYMMDDHHIISST)]
	} else if len(data) >= len(Layout_YYYYMMDDHHIISS) {
		layout = Layout_YYYYMMDDHHIISS
		data = data[:len(Layout_YYYYMMDDHHIISS)]
	} else if len(data) >= len(Layout_YYYYMMDD) {
		layout = Layout_YYYYMMDD
		data = data[:len(Layout_YYYYMMDD)]
	} else if len(data) >= len(Layout_YYYYMM) {
		layout = Layout_YYYYMM
		data = data[:len(Layout_YYYYMM)]
	}
	now, _ := time.ParseInLocation(layout, string(data), time.Local)
	if now.UnixNano() < 0 || now.Unix() <= 0 {
		*t = Time{}
	} else {
		*t = Time{now}
	}
	return nil
}

// deserialization.
func (d *Time) UnmarshalText(text []byte) error {
	*d = ToTime(string(text))
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for XML
// serialization.
func (d Time) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

// Value insert timestamp into mysql need this function.
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *Time) Scan(v interface{}) error {
	*t = ToTime(v)
	return nil
}
func (d Time) Valid() bool {
	return !d.IsZero()
}
