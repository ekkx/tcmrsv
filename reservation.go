package tcmrsv

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func (rsv *TCMRSV) GetMyReservations() ([]Reservation, error) {
	req, err := http.NewRequest(http.MethodGet, ENDPOINT_INDEX, nil)
	if err != nil {
		return nil, err
	}

	res, err := rsv.DoRequest(req)
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
						currentReservation.Date = text
					case "res-time":
						currentReservation.TimeRange = text
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
	Date       time.Time
	FromHour   int
	FromMinute int
	ToHour     int
	ToMinute   int
}

func (rsv *TCMRSV) Reserve(params *ReserveParams) (*http.Response, error) {
	u, err := url.Parse(ENDPOINT_CONFIRMS)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("campus", string(params.Campus))
	q.Set("room", params.RoomID)
	q.Set("ymd", params.Date.Format("2006/01/02 15:04:05"))
	q.Set("fromh", fmt.Sprintf("%02d", params.FromHour))
	q.Set("fromm", fmt.Sprintf("%02d", params.FromMinute)) // TODO: 00 か 30 バリデートする
	q.Set("toh", fmt.Sprintf("%02d", params.ToHour))
	q.Set("tom", fmt.Sprintf("%02d", params.ToMinute)) // TODO: 00 か 30 バリデートする
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := rsv.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	rsv.aspcfg.Update(res.Body)

	form := url.Values{}
	form.Set("__VIEWSTATE", rsv.aspcfg.ViewState)
	form.Set("__VIEWSTATEGENERATOR", rsv.aspcfg.ViewStateGenerator)
	form.Set("__EVENTVALIDATION", rsv.aspcfg.EventValidation)
	form.Set("KakuteiButton", "")

	req, err = http.NewRequest(http.MethodPost, u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return rsv.DoRequest(req)
}

type CancelReservationParams struct {
	ReservationID string
	Comment       string
}

func (rsv *TCMRSV) CancelReservation(params *CancelReservationParams) (*http.Response, error) {
	u, err := url.Parse(ENDPOINT_CANCEL_RESERVATION)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("id", string(params.ReservationID))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := rsv.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	rsv.aspcfg.Update(res.Body)

	form := url.Values{}
	form.Set("__EVENTTARGET", "")
	form.Set("__EVENTARGUMENT", "")
	form.Set("__VIEWSTATE", rsv.aspcfg.ViewState)
	form.Set("__VIEWSTATEGENERATOR", rsv.aspcfg.ViewStateGenerator)
	form.Set("__EVENTVALIDATION", rsv.aspcfg.EventValidation)
	form.Set("freeword", params.Comment)
	form.Set("YoyakuCancelButton", "")

	req, err = http.NewRequest(http.MethodPost, u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return rsv.DoRequest(req)
}
