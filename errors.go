package tcmrsv

import (
	"errors"
	"strings"

	"io"

	"golang.org/x/net/html"
)

var (
	ErrAuthenticationFailed    = errors.New("authentication failed error")
	ErrCreateReservationFailed = errors.New("create reservation failed error")
	ErrCancelReservationFailed = errors.New("cancel reservation failed error")
	ErrInvalidCampus           = errors.New("invalid campus error")
	ErrInvalidIDFormat         = errors.New("invalid ID format error")
	ErrDateOutOfRange          = errors.New("date out of range error")
	ErrInvalidTimeRange        = errors.New("invalid time range error")
	ErrTimeInPast              = errors.New("time in past error")
	ErrInvalidComment          = errors.New("invalid comment error")
	ErrInternalServer          = errors.New("internal server error")
)

func isInternalServerErrorPage(body io.Reader) (bool, error) {
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

func isLoginPage(body io.Reader) (bool, error) {
	z := html.NewTokenizer(body)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return false, nil
			}
			return false, z.Err()

		case html.SelfClosingTagToken, html.StartTagToken:
			t := z.Token()
			if t.Data == "input" {
				for _, attr := range t.Attr {
					if attr.Key == "id" && attr.Val == "btnLogin" {
						return true, nil
					}
				}
			}
		}
	}
}
