package tcmrsv

import (
	"testing"
	"time"
)

func TestNewDate(t *testing.T) {
	tests := []struct {
		name  string
		year  int
		month time.Month
		day   int
		want  Date
	}{
		{
			name:  "valid date",
			year:  2024,
			month: time.January,
			day:   15,
			want:  Date{Year: 2024, Month: time.January, Day: 15},
		},
		{
			name:  "leap year date",
			year:  2024,
			month: time.February,
			day:   29,
			want:  Date{Year: 2024, Month: time.February, Day: 29},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDate(tt.year, tt.month, tt.day)
			if got != tt.want {
				t.Errorf("NewDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_IsZero(t *testing.T) {
	tests := []struct {
		name string
		date Date
		want bool
	}{
		{
			name: "zero date",
			date: Date{},
			want: true,
		},
		{
			name: "non-zero date",
			date: NewDate(2024, time.January, 1),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date.IsZero(); got != tt.want {
				t.Errorf("Date.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromTime(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	tests := []struct {
		name string
		time time.Time
		want Date
	}{
		{
			name: "JST time",
			time: time.Date(2024, time.March, 15, 10, 30, 45, 0, jst),
			want: Date{Year: 2024, Month: time.March, Day: 15},
		},
		{
			name: "UTC time converted to JST",
			time: time.Date(2024, time.March, 15, 0, 0, 0, 0, time.UTC),
			want: Date{Year: 2024, Month: time.March, Day: 15},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromTime(tt.time)
			if got != tt.want {
				t.Errorf("FromTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_ToTime(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	tests := []struct {
		name string
		date Date
		want time.Time
	}{
		{
			name: "regular date",
			date: NewDate(2024, time.March, 15),
			want: time.Date(2024, time.March, 15, 0, 0, 0, 0, jst),
		},
		{
			name: "leap year date",
			date: NewDate(2024, time.February, 29),
			want: time.Date(2024, time.February, 29, 0, 0, 0, 0, jst),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.date.ToTime()
			if !got.Equal(tt.want) {
				t.Errorf("Date.ToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_String(t *testing.T) {
	tests := []struct {
		name string
		date Date
		want string
	}{
		{
			name: "single digit month and day",
			date: NewDate(2024, time.January, 5),
			want: "2024-01-05",
		},
		{
			name: "double digit month and day",
			date: NewDate(2024, time.December, 31),
			want: "2024-12-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date.String(); got != tt.want {
				t.Errorf("Date.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Date
		wantErr bool
	}{
		{
			name:    "valid date",
			input:   "2024-03-15",
			want:    NewDate(2024, time.March, 15),
			wantErr: false,
		},
		{
			name:    "invalid format - wrong separator",
			input:   "2024/03/15",
			want:    Date{},
			wantErr: true,
		},
		{
			name:    "invalid format - wrong order",
			input:   "15-03-2024",
			want:    Date{},
			wantErr: true,
		},
		{
			name:    "invalid date",
			input:   "2024-13-01",
			want:    Date{},
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    Date{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDate() = %v, want %v", got, tt.want)
			}
			if tt.wantErr && err != ErrInvalidDateFormat {
				t.Errorf("ParseDate() error = %v, want %v", err, ErrInvalidDateFormat)
			}
		})
	}
}

func TestDate_IsValid(t *testing.T) {
	tests := []struct {
		name string
		date Date
		want bool
	}{
		{
			name: "valid date",
			date: NewDate(2024, time.March, 15),
			want: true,
		},
		{
			name: "invalid date - February 30",
			date: Date{Year: 2024, Month: time.February, Day: 30},
			want: false,
		},
		{
			name: "valid leap year date",
			date: NewDate(2024, time.February, 29),
			want: true,
		},
		{
			name: "invalid non-leap year date",
			date: Date{Year: 2023, Month: time.February, Day: 29},
			want: false,
		},
		{
			name: "invalid month",
			date: Date{Year: 2024, Month: 13, Day: 1},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date.IsValid(); got != tt.want {
				t.Errorf("Date.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Equals(t *testing.T) {
	tests := []struct {
		name  string
		date1 Date
		date2 Date
		want  bool
	}{
		{
			name:  "equal dates",
			date1: NewDate(2024, time.March, 15),
			date2: NewDate(2024, time.March, 15),
			want:  true,
		},
		{
			name:  "different dates",
			date1: NewDate(2024, time.March, 15),
			date2: NewDate(2024, time.March, 16),
			want:  false,
		},
		{
			name:  "different months",
			date1: NewDate(2024, time.March, 15),
			date2: NewDate(2024, time.April, 15),
			want:  false,
		},
		{
			name:  "different years",
			date1: NewDate(2024, time.March, 15),
			date2: NewDate(2025, time.March, 15),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date1.Equals(tt.date2); got != tt.want {
				t.Errorf("Date.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_IsBefore(t *testing.T) {
	tests := []struct {
		name  string
		date1 Date
		date2 Date
		want  bool
	}{
		{
			name:  "date1 before date2",
			date1: NewDate(2024, time.March, 14),
			date2: NewDate(2024, time.March, 15),
			want:  true,
		},
		{
			name:  "date1 after date2",
			date1: NewDate(2024, time.March, 16),
			date2: NewDate(2024, time.March, 15),
			want:  false,
		},
		{
			name:  "equal dates",
			date1: NewDate(2024, time.March, 15),
			date2: NewDate(2024, time.March, 15),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date1.IsBefore(tt.date2); got != tt.want {
				t.Errorf("Date.IsBefore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_IsAfter(t *testing.T) {
	tests := []struct {
		name  string
		date1 Date
		date2 Date
		want  bool
	}{
		{
			name:  "date1 after date2",
			date1: NewDate(2024, time.March, 16),
			date2: NewDate(2024, time.March, 15),
			want:  true,
		},
		{
			name:  "date1 before date2",
			date1: NewDate(2024, time.March, 14),
			date2: NewDate(2024, time.March, 15),
			want:  false,
		},
		{
			name:  "equal dates",
			date1: NewDate(2024, time.March, 15),
			date2: NewDate(2024, time.March, 15),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date1.IsAfter(tt.date2); got != tt.want {
				t.Errorf("Date.IsAfter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToday(t *testing.T) {
	today := Today()
	now := time.Now()
	jst, _ := time.LoadLocation("Asia/Tokyo")
	nowJST := now.In(jst)

	expectedDate := NewDate(nowJST.Year(), nowJST.Month(), nowJST.Day())

	if !today.Equals(expectedDate) {
		t.Errorf("Today() = %v, want %v", today, expectedDate)
	}
}

func TestDate_AddDays(t *testing.T) {
	tests := []struct {
		name string
		date Date
		days int
		want Date
	}{
		{
			name: "add positive days",
			date: NewDate(2024, time.March, 15),
			days: 5,
			want: NewDate(2024, time.March, 20),
		},
		{
			name: "add negative days",
			date: NewDate(2024, time.March, 15),
			days: -5,
			want: NewDate(2024, time.March, 10),
		},
		{
			name: "add days crossing month boundary",
			date: NewDate(2024, time.March, 30),
			days: 3,
			want: NewDate(2024, time.April, 2),
		},
		{
			name: "add days crossing year boundary",
			date: NewDate(2024, time.December, 30),
			days: 3,
			want: NewDate(2025, time.January, 2),
		},
		{
			name: "add zero days",
			date: NewDate(2024, time.March, 15),
			days: 0,
			want: NewDate(2024, time.March, 15),
		},
		{
			name: "leap year boundary",
			date: NewDate(2024, time.February, 28),
			days: 1,
			want: NewDate(2024, time.February, 29),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.date.AddDays(tt.days)
			if !got.Equals(tt.want) {
				t.Errorf("Date.AddDays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustLoadJST(t *testing.T) {
	jst := mustLoadJST()
	if jst == nil {
		t.Error("mustLoadJST() returned nil")
		return
	}

	name, offset := time.Now().In(jst).Zone()
	expectedOffset := 9 * 60 * 60

	if offset != expectedOffset {
		t.Errorf("JST offset = %d, want %d", offset, expectedOffset)
	}

	if name != "JST" && name != "Asia/Tokyo" {
		t.Errorf("JST zone name = %s, want JST or Asia/Tokyo", name)
	}
}
