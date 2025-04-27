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
		},
		aspcfg: NewASPConfig(),
	}
}
