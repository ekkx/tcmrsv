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

func (c *Client) Login(params *LoginParams) error {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+ENDPOINT_LOGIN, nil)
	if err != nil {
		return err
	}

	res, err := c.DoRequest(req, false)
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
	form.Set("input_id", params.UserID)
	form.Set("input_pass", params.Password)
	form.Set("btnLogin", "")

	req, err = http.NewRequest(http.MethodPost, c.baseURL+ENDPOINT_LOGIN, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err = c.DoRequest(req, true)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
