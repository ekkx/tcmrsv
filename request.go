package tcmrsv

import (
	"bytes"
	"io"
	"net/http"
)

func (rsv *TCMRSV) DoRequest(req *http.Request, requireAuth bool) (*http.Response, error) {
	res, err := rsv.client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	reader := func() *bytes.Reader {
		return bytes.NewReader(bodyBytes)
	}

	isErr, err := isInternalServerErrorPage(reader())
	if err != nil {
		return nil, err
	}
	if isErr {
		return nil, ErrInternalServer
	}

	if requireAuth {
		isAuthErr, err := isLoginPage(reader())
		if err != nil {
			return nil, err
		}
		if isAuthErr {
			return nil, ErrAuthenticationFailed
		}
	}

	if err := rsv.aspcfg.Update(reader()); err != nil {
		return nil, err
	}

	res.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	return res, nil
}
