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
	Date   Date
}

// 利用可能な練習室一覧を取得する
func (c *Client) GetRoomAvailability(params *GetRoomAvailabilityParams) ([]RoomAvailability, error) {
	now := time.Now().In(jst)

	if !params.Campus.IsValid() {
		return nil, ErrInvalidCampus
	}
	if !IsDateWithin2Days(now, params.Date) {
		return nil, ErrInvalidTimeRange
	}

	u, err := url.Parse(c.baseURL + ENDPOINT_RESERVE)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("campus", string(params.Campus))
	q.Set("ymd", params.Date.ToTime().Format("2006/01/02 15:04:05"))
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

	isToday := params.Date.Equals(Today())

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

				class := attrs["class"]
				isJudgment4 := strings.HasPrefix(class, "judgment4")

				// 部屋
				if colIndex == 1 {
					continue
				}

				// 空き時間
				if colIndex >= 2 && currentAvailability != nil && isJudgment4 {
					tdHasAvailableInput := false
					tdHasCircle := false
					depth := 1

					for depth > 0 {
						tt := z.Next()
						switch tt {
						case html.ErrorToken:
							depth = 0
						case html.StartTagToken, html.SelfClosingTagToken:
							t := z.Token()
							if t.Data == "input" {
								disabled := false
								for _, a := range t.Attr {
									if a.Key == "disabled" {
										disabled = true
										break
									}
								}
								if !disabled {
									tdHasAvailableInput = true
								}
							}
						case html.TextToken:
							text := strings.TrimSpace(z.Token().Data)
							if strings.Contains(text, "〇") {
								tdHasCircle = true
							}
						case html.EndTagToken:
							t := z.Token()
							if t.Data == "td" {
								depth--
							}
						}
					}

					if tdHasAvailableInput || tdHasCircle {
						offset := colIndex - 2
						hour := baseHour + (offset*30)/60
						minute := (offset * 30) % 60

						if isToday {
							slotTime := time.Date(params.Date.Year, params.Date.Month, params.Date.Day, hour, minute, 0, 0, jst)
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
						for _, r := range c.GetRooms() {
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
