package tcmrsv

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
)

type MockServer struct {
	Server   *httptest.Server
	Client   *Client
	Requests []*http.Request // リクエスト検証用
}

func NewMockServer(handler http.HandlerFunc) *MockServer {
	ms := &MockServer{
		Requests: make([]*http.Request, 0),
	}

	recordingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCopy := &http.Request{
			Method: r.Method,
			URL:    r.URL,
			Header: r.Header.Clone(),
		}

		if r.Body != nil {
			bodyBytes, _ := io.ReadAll(r.Body)
			r.Body.Close()
			reqCopy.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		ms.Requests = append(ms.Requests, reqCopy)

		handler(w, r)
	})

	ms.Server = httptest.NewServer(recordingHandler)

	ms.Client = New(
		WithBaseURL(ms.Server.URL),
		WithHTTPClient(ms.Server.Client()),
	)

	return ms
}

func (m *MockServer) Close() {
	m.Server.Close()
}

func LoadFixture(path string) string {
	data, err := os.ReadFile(filepath.Join("tests", "fixtures", path))
	if err != nil {
		panic(err)
	}
	return string(data)
}

func CreateHandler(routes map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path
		if handler, ok := routes[key]; ok {
			handler(w, r)
			return
		}

		key = r.Method + " " + r.URL.Path
		if handler, ok := routes[key]; ok {
			handler(w, r)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}
}
