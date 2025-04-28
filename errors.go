package tcmrsv

import (
	"errors"
	"strings"

	"io"

	"golang.org/x/net/html"
)

var (
	ErrInternalServer = errors.New("internal server error")
)

func isErrorPage(body io.Reader) (bool, error) {
	z := html.NewTokenizer(body)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return false, nil
			}
			return false, z.Err()

		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "span" {
				for _, attr := range t.Attr {
					if attr.Key == "class" && attr.Val == "title" {
						tt = z.Next()
						if tt == html.TextToken {
							text := strings.TrimSpace(string(z.Text()))
							if strings.Contains(text, "現在、サーバへのアクセスが集中し、ページを閲覧しにくい状態になっております。") {
								return true, nil
							}
						}
					}
				}
			}
		}
	}
}
