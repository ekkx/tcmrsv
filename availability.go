package tcmrsv

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type GetRoomAvailabilityParams struct {
	Campus Campus
	Date   time.Time
}

// 利用可能な練習室一覧を取得する
func (c *Client) GetRoomAvailability(params *GetRoomAvailabilityParams) ([]RoomAvailability, error) {
	if !params.Campus.IsValid() {
		return nil, ErrInvalidCampus
	}
	if IsDateWithin2Days(params.Date) {
		return nil, ErrInvalidTimeRange
	}

	u, err := url.Parse(c.baseURL+ENDPOINT_RESERVE)
	if err != nil {
		return nil, err
	}

	jstDate := time.Date(params.Date.Year(), params.Date.Month(), params.Date.Day(), 0, 0, 0, 0, jst())

	q := u.Query()
	q.Set("campus", string(params.Campus))
	q.Set("ymd", jstDate.Format("2006/01/02 15:04:05"))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req, true)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	z := html.NewTokenizer(res.Body)

	var availabilities []RoomAvailability
	var insideTable, insideRow bool
	var colIndex int
	var currentRoomName string
	var currentAvailability *RoomAvailability

	const baseHour = 7

	now := time.Now().In(jst())
	isToday := params.Date.Year() == now.Year() &&
		params.Date.Month() == now.Month() &&
		params.Date.Day() == now.Day()

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return availabilities, nil
			}
			return nil, z.Err()

		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()

			switch t.Data {
			case "table":
				for _, attr := range t.Attr {
					if attr.Key == "id" && strings.HasPrefix(attr.Val, "aspTable") {
						insideTable = true
					}
				}

			case "tr":
				if insideTable {
					colIndex = 0
					insideRow = true
					currentRoomName = ""
					currentAvailability = nil
				}

			case "td":
				if !insideRow {
					break
				}
				colIndex++
				attrs := map[string]string{}
				for _, a := range t.Attr {
					attrs[a.Key] = a.Val
				}

				// 部屋
				if colIndex == 1 {
					continue
				}

				// 空き時間
				if colIndex >= 2 && currentAvailability != nil {
					class := attrs["class"]
					style := attrs["style"]
					isAvailable := strings.HasPrefix(class, "judgment4") || style == "text-align:center;"
					disabled := false

					if isAvailable && !disabled {
						offset := colIndex - 2
						hour := baseHour + (offset*30)/60
						minute := (offset * 30) % 60

						if isToday {
							slotTime := time.Date(params.Date.Year(), params.Date.Month(), params.Date.Day(), hour, minute, 0, 0, jst())
							if slotTime.Before(now) {
								break
							}
						}

						currentAvailability.AvailableTimes = append(currentAvailability.AvailableTimes, AvailableTime{
							Hour:   hour,
							Minute: minute,
						})
					}
				}

			case "span":
				if insideRow && colIndex == 1 {
					if z.Next() == html.TextToken {
						currentRoomName = strings.TrimSpace(z.Token().Data)
						for _, r := range GetRooms() {
							if r.Name == currentRoomName {
								currentAvailability = &RoomAvailability{Room: r}
								break
							}
						}
					}
				}
			}

		case html.EndTagToken:
			t := z.Token()
			switch t.Data {
			case "tr":
				if insideRow {
					if currentAvailability != nil && len(currentAvailability.AvailableTimes) > 0 {
						availabilities = append(availabilities, *currentAvailability)
					}
					insideRow = false
				}
			case "table":
				insideTable = false
			}
		}
	}
}
