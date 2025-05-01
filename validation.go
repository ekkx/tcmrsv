package tcmrsv

import (
	"regexp"
	"strings"
	"time"
)

func IsIDValid(ID string) bool {
	idRegex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	return idRegex.MatchString(ID)
}

func IsDateWithin2Days(date time.Time) bool {
	now := time.Now().In(jst())
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst())
	dayAfterTomorrow := today.AddDate(0, 0, 3)

	return date.After(today) && date.Before(dayAfterTomorrow)
}

func IsTimeRangeValid(fromHour, fromMinute, toHour, toMinute int) bool {
	if fromHour < 7 || fromHour > 22 {
		return false
	}

	if toHour < fromHour || toHour > 23 {
		return false
	}

	if fromMinute != 0 && fromMinute != 30 {
		return false
	}

	if toMinute != 0 && toMinute != 30 {
		return false
	}

	if toHour == 23 && toMinute != 0 {
		return false
	}

	fromTotal := fromHour*60 + fromMinute
	toTotal := toHour*60 + toMinute

	return toTotal > fromTotal
}

func IsTimeInFuture(fromHour, fromMinute int) bool {
	now := time.Now().In(jst())
	currentTotal := now.Hour()*60 + now.Minute()
	fromTotal := fromHour*60 + fromMinute

	return fromTotal >= currentTotal
}

func IsCommentValid(comment string) bool {
	return strings.TrimSpace(comment) != ""
}
