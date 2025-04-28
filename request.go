package tcmrsv

import (
	"bytes"
	"io"
	"net/http"
)

func (rsv *TCMRSV) DoRequest(req *http.Request) (*http.Response, error) {
	res, err := rsv.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil && res.Body != nil {
			res.Body.Close()
		}
	}()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	isErr, err := isErrorPage(bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	if isErr {
		return nil, ErrInternalServer
	}

	res.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	return res, nil
}
