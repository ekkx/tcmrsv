package tcmrsv

import (
	"net/http"
	"net/url"
)

type LoginParams struct {
	UserID   string
	Password string
}

func (rsv *TCMRSV) Login(params *LoginParams) (*http.Response, error) {
	res, err := rsv.client.Get(ENDPOINT_LOGIN)
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
	form.Set("input_id", params.UserID)
	form.Set("input_pass", params.Password)
	form.Set("btnLogin", "")

	res, err = rsv.client.PostForm(ENDPOINT_LOGIN, form)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return res, nil
}
