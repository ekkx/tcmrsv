package tcmrsv

import (
	"net/http"
	"net/http/cookiejar"
)

type TCMRSV struct {
	client *http.Client
	aspcfg *ASPConfig
}

func New() *TCMRSV {
	jar, _ := cookiejar.New(nil)

	return &TCMRSV{
		client: &http.Client{
			Jar: jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				req.URL.RawQuery = req.URL.Query().Encode()
				return nil
			},
		},
		aspcfg: NewASPConfig(),
	}
}
