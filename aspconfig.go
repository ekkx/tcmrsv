package tcmrsv

import (
	"io"

	"golang.org/x/net/html"
)

type ASPConfig struct {
	ViewState          string
	ViewStateGenerator string
	EventValidation    string
}

func NewASPConfig() *ASPConfig {
	return &ASPConfig{}
}

func (cfg *ASPConfig) Update(r io.Reader) error {
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return nil
			}
			return z.Err()

		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()

			if t.Data == "input" {
				attrs := make(map[string]string)
				for _, a := range t.Attr {
					attrs[a.Key] = a.Val
				}

				switch {
				case attrs["id"] == "__VIEWSTATE":
					cfg.ViewState = attrs["value"]
				case attrs["id"] == "__VIEWSTATEGENERATOR":
					cfg.ViewStateGenerator = attrs["value"]
				case attrs["id"] == "__EVENTVALIDATION":
					cfg.EventValidation = attrs["value"]
				}
			}
		}
	}
}
