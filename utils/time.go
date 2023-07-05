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
func _stringToTimeString(timestr string) string {
	var ret = ""
	timestr = strings.Replace(timestr, `"`, "", -1)
	timestr = strings.Replace(timestr, "/", "-", -1)
	//2022-01-26T16:33:18.000+0800
	if len(timestr) >= len(Layout_YYYYMMDDHHIISST) && strings.Contains(timestr, "T") {
		timestr = strings.Replace(timestr, "T", " ", -1)
		timestr = timestr[:len(Layout_YYYYMMDDHHIISST)]
	}
	arr := strings.Split(timestr, " ")
	if len(arr) <= 1 {
		ret = strings.Join([]string{arr[0], "00:00:00"}, " ")
	} else {
		switch strings.Count(arr[1], ":") {
		case 0:
			ret = strings.Join([]string{arr[0], strings.Join([]string{arr[1], ":00:00"}, "")}, " ")
			break
		case 1:
			ret = strings.Join([]string{arr[0], strings.Join([]string{arr[1], ":00"}, "")}, " ")
			break
		default:
			ret = timestr
			break
		}
	}
	return ret
}
func ToTime(value interface{}) Time {
	if value == "" || value == "null" || value == nil {
		return Time{}
	}
	if v, ok := value.(time.Time); ok {
		return Time{v}
	}
	if v, ok := value.(Time); ok {
		return v
	}
	if v, ok := value.(string); ok {
		v = _stringToTimeString(v)
		layout := "2006-1-2 15:4:5"
		data := []rune(v)
		now, err := time.ParseInLocation(layout, string(data), time.Local)
		if err != nil {
			log.Print("string to ToTime failed", err)
		}
		return Time{now}
	}
	return Time{}
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
	*t = ToTime(string(data))
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
