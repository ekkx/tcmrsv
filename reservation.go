package tcmrsv

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Campus int

const (
	CampusIkebukuro  Campus = 1
	CampusNakameguro Campus = 2
)

type ReserveParams struct {
	Campus     Campus
	RoomID     string
	Start      time.Time
	FromHour   int
	FromMinute int
	ToHour     int
	ToMinute   int
}

func (rsv *TCMRSV) Reserve(params *ReserveParams) (*http.Response, error) {
	u, err := url.Parse(ENDPOINT_RESERVATION)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	q.Set("campus", fmt.Sprintf("%d", params.Campus))
	q.Set("room", params.RoomID)
	q.Set("ymd", params.Start.Format("2006/01/02 15:04:05"))
	q.Set("fromh", fmt.Sprintf("%02d", params.FromHour))
	q.Set("fromm", fmt.Sprintf("%02d", params.FromMinute)) // TODO: 00 か 30 バリデートする
	q.Set("toh", fmt.Sprintf("%02d", params.ToHour))
	q.Set("tom", fmt.Sprintf("%02d", params.ToMinute)) // TODO: 00 か 30 バリデートする
	u.RawQuery = q.Encode()

	confirmUrl := u.String()
	res, err := rsv.client.Get(confirmUrl)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	rsv.aspcfg.Update(res.Body)

	form := url.Values{}
	form.Set("__VIEWSTATE", rsv.aspcfg.ViewState)
	form.Set("__VIEWSTATEGENERATOR", rsv.aspcfg.ViewStateGenerator)
	form.Set("__EVENTVALIDATION", rsv.aspcfg.EventValidation)
	form.Set("KakuteiButton", "")

	res, err = rsv.client.PostForm(confirmUrl, form)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return res, nil
}
