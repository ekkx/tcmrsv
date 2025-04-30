package tcmrsv

import (
	"net/http"
	"net/url"
	"strings"
)

type LoginParams struct {
	UserID   string
	Password string
}

func (rsv *TCMRSV) Login(params *LoginParams) error {
	req, err := http.NewRequest(http.MethodGet, ENDPOINT_LOGIN, nil)
	if err != nil {
		return err
	}

	res, err := rsv.DoRequest(req, false)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	form := url.Values{}
	form.Set("__EVENTTARGET", "")
	form.Set("__EVENTARGUMENT", "")
	form.Set("__VIEWSTATE", rsv.aspcfg.ViewState)
	form.Set("__VIEWSTATEGENERATOR", rsv.aspcfg.ViewStateGenerator)
	form.Set("__EVENTVALIDATION", rsv.aspcfg.EventValidation)
	form.Set("input_id", params.UserID)
	form.Set("input_pass", params.Password)
	form.Set("btnLogin", "")

	req, err = http.NewRequest(http.MethodPost, ENDPOINT_LOGIN, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err = rsv.DoRequest(req, true)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
