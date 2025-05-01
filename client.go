package tcmrsv

import (
	"bytes"
	"io"
	"net/http"
	"net/http/cookiejar"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	aspConfig  *ASPConfig
}

type ClientConfig struct {
	httpClient *http.Client
	baseURL    string
}

func newClientConfig() *ClientConfig {
	jar, _ := cookiejar.New(nil)

	return &ClientConfig{
		httpClient: &http.Client{
			Jar: jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				req.URL.RawQuery = req.URL.Query().Encode()
				return nil
			},
		},
		baseURL: "https://www.tokyo-ondai-career.jp",
	}
}

type ClientOption func(cfg *ClientConfig)

func WithBaseURL(baseURL string) ClientOption {
	return func(cfg *ClientConfig) {
		cfg.baseURL = baseURL
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(cfg *ClientConfig) {
		if httpClient != nil {
			cfg.httpClient = httpClient
		}
	}
}

func New(options ...ClientOption) *Client {
	cfg := newClientConfig()
	for _, opt := range options {
		opt(cfg)
	}

	return &Client{
		httpClient: cfg.httpClient,
		baseURL:    cfg.baseURL,
		aspConfig:  NewASPConfig(),
	}
}

func (c *Client) DoRequest(req *http.Request, requireAuth bool) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
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

	if err := c.aspConfig.Update(reader()); err != nil {
		return nil, err
	}

	res.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	return res, nil
}
