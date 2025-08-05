package tcmrsv

import (
	"errors"
	"time"
)

const Layout = "2006-01-02"

var (
	ErrInvalidDateFormat = errors.New("invalid date format, expected YYYY-MM-DD")
)

var jst = mustLoadJST()

func mustLoadJST() *time.Location {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		loc = time.FixedZone("JST", 9*60*60)
	}
	return loc
}

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func NewDate(year int, month time.Month, day int) Date {
	return Date{Year: year, Month: month, Day: day}
}

func (d Date) IsZero() bool {
	return d == Date{}
}

func FromTime(t time.Time) Date {
	t = t.In(jst)
	return NewDate(t.Year(), t.Month(), t.Day())
}

func (d Date) ToTime() time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, jst)
}

func (d Date) String() string {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, jst).Format(Layout)
}

func ParseDate(s string) (Date, error) {
	t, err := time.ParseInLocation(Layout, s, jst)
	if err != nil {
		return Date{}, ErrInvalidDateFormat
	}
	return FromTime(t), nil
}

func (d Date) IsValid() bool {
	t := time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, jst)
	return t.Year() == d.Year && t.Month() == d.Month && t.Day() == d.Day
}

func (d Date) Equals(other Date) bool {
	return d.ToTime().Equal(other.ToTime())
}

func (d Date) IsBefore(other Date) bool {
	return d.ToTime().Before(other.ToTime())
}

func (d Date) IsAfter(other Date) bool {
	return d.ToTime().After(other.ToTime())
}

func Today() Date {
	return FromTime(time.Now().In(jst))
}

func (d Date) AddDays(days int) Date {
	t := d.ToTime().AddDate(0, 0, days)
	return FromTime(t)
}
