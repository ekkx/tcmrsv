package tcmrsv

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func (c *Client) GetMyReservations() ([]Reservation, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+ENDPOINT_INDEX, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req, true)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	z := html.NewTokenizer(res.Body)

	var (
		reservations          []Reservation
		insideReservationList bool
		currentReservation    *Reservation
		currentTag            string
	)

	var campusMap = map[string]Campus{
		"池袋キャンパス":      CampusIkebukuro,
		"中目黒・代官山キャンパス": CampusNakameguro,
	}

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return reservations, nil
			}
			return nil, z.Err()

		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()

			switch t.Data {
			case "div":
				for _, attr := range t.Attr {
					if attr.Key == "id" && attr.Val == "reservation-list" {
						insideReservationList = true
					}
				}

			case "dl":
				if insideReservationList {
					currentReservation = &Reservation{}
				}

			case "dt", "dd":
				currentTag = ""
				for _, attr := range t.Attr {
					switch attr.Key {
					case "class":
						if t.Data == "dt" && attr.Val == "res-room" {
							currentTag = "campus"
						} else {
							currentTag = attr.Val
						}
					case "style":
						if strings.Contains(attr.Val, "width:160px") {
							currentTag = "campus"
						}
					}
				}
				if currentTag == "" {
					currentTag = "campus"
				}

			case "a":
				if insideReservationList && currentReservation != nil && currentTag == "res-cancell" {
					for _, attr := range t.Attr {
						if attr.Key == "href" {
							if idx := strings.Index(attr.Val, "id="); idx >= 0 {
								id := attr.Val[idx+3:]
								currentReservation.ID = strings.TrimSpace(id)
							}
						}
					}
				}
			}

		case html.TextToken:
			if insideReservationList && currentReservation != nil {
				text := strings.TrimSpace(z.Token().Data)
				if text != "" {
					switch currentTag {
					case "campus":
						if campus, ok := campusMap[text]; ok {
							currentReservation.Campus = campus
						} else {
							currentReservation.Campus = CampusUnknown
						}
						currentReservation.CampusName = text
					case "res-date":
						layout := "2006年01月02日"
						text = strings.SplitN(text, "（", 2)[0]
						parsedDate, err := time.Parse(layout, text)
						if err != nil {
							currentReservation.Date = Date{} // Fallback to raw text if parsing fails
						} else {
							currentReservation.Date = FromTime(parsedDate)
						}
					case "res-time":
						times := strings.Split(text, "-")
						if len(times) != 2 {
							continue // Invalid time format
						}
						fromTime := strings.Split(times[0], ":")
						toTime := strings.Split(times[1], ":")
						if len(fromTime) != 2 || len(toTime) != 2 {
							continue // Invalid time format
						}
						fromHour, err1 := strconv.Atoi(fromTime[0])
						fromMinute, err2 := strconv.Atoi(fromTime[1])
						toHour, err3 := strconv.Atoi(toTime[0])
						toMinute, err4 := strconv.Atoi(toTime[1])
						if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
							continue // Invalid time format
						}
						currentReservation.FromHour = fromHour
						currentReservation.FromMinute = fromMinute
						currentReservation.ToHour = toHour
						currentReservation.ToMinute = toMinute
					case "res-room":
						currentReservation.RoomName = text
					}
				}
			}

		case html.EndTagToken:
			t := z.Token()
			if t.Data == "dl" && insideReservationList && currentReservation != nil {
				// 変な空構造体が入るので簡単に検証
				if currentReservation.CampusName != "" || currentReservation.ID != "" {
					reservations = append(reservations, *currentReservation)
				}
				currentReservation = nil
			}
		}
	}
}

type ReserveParams struct {
	Campus     Campus
	RoomID     string
	Date       Date
	FromHour   int
	FromMinute int
	ToHour     int
	ToMinute   int
}

func (c *Client) Reserve(params *ReserveParams) error {
	if !params.Campus.IsValid() {
		return ErrInvalidCampus
	}
	if !IsIDValid(params.RoomID) {
		return ErrInvalidIDFormat
	}
	if !IsDateWithin2Days(time.Now().In(jst), params.Date) {
		return ErrDateOutOfRange
	}
	if !IsTimeRangeValid(params.FromHour, params.FromMinute, params.ToHour, params.ToMinute) {
		return ErrInvalidTimeRange
	}
	if !IsTimeInFuture(params.FromHour, params.FromMinute, params.Date) {
		return ErrTimeInPast
	}

	u, err := url.Parse(c.baseURL + ENDPOINT_CONFIRMS)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("campus", string(params.Campus))
	q.Set("room", params.RoomID)
	q.Set("ymd", params.Date.ToTime().Format("2006/01/02 15:04:05"))
	q.Set("fromh", fmt.Sprintf("%02d", params.FromHour))
	q.Set("fromm", fmt.Sprintf("%02d", params.FromMinute))
	q.Set("toh", fmt.Sprintf("%02d", params.ToHour))
	q.Set("tom", fmt.Sprintf("%02d", params.ToMinute))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	res, err := c.DoRequest(req, true)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	form := url.Values{}
	form.Set("__VIEWSTATE", c.aspConfig.ViewState)
	form.Set("__VIEWSTATEGENERATOR", c.aspConfig.ViewStateGenerator)
	form.Set("__EVENTVALIDATION", c.aspConfig.EventValidation)
	form.Set("KakuteiButton", "")

	req, err = http.NewRequest(http.MethodPost, u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err = c.DoRequest(req, true)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if !strings.Contains(string(bodyBytes), "予約が完了しました") {
		return ErrCreateReservationFailed
	}

	return nil
}

type CancelReservationParams struct {
	ReservationID string
	Comment       string
}

func (c *Client) CancelReservation(params *CancelReservationParams) error {
	if !IsIDValid(params.ReservationID) {
		return ErrInvalidIDFormat
	}
	if !IsCommentValid(params.Comment) {
		return ErrInvalidComment
	}

	u, err := url.Parse(c.baseURL + ENDPOINT_CANCEL_RESERVATION)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("id", string(params.ReservationID))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	res, err := c.DoRequest(req, true)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	form := url.Values{}
	form.Set("__EVENTTARGET", "")
	form.Set("__EVENTARGUMENT", "")
	form.Set("__VIEWSTATE", c.aspConfig.ViewState)
	form.Set("__VIEWSTATEGENERATOR", c.aspConfig.ViewStateGenerator)
	form.Set("__EVENTVALIDATION", c.aspConfig.EventValidation)
	form.Set("freeword", params.Comment)
	form.Set("YoyakuCancelButton", "")

	req, err = http.NewRequest(http.MethodPost, u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err = c.DoRequest(req, true)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if !strings.Contains(string(bodyBytes), "予約キャンセル完了") {
		return ErrCancelReservationFailed
	}

	return nil
}
