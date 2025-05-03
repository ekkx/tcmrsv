package tcmrsv

import (
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Run("SuccessfulLogin", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /index.aspx": func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(LoadFixture("index.html")))
			},
			"POST /index.aspx": func(w http.ResponseWriter, r *http.Request) {
				// ログイン成功後は予約一覧ページにリダイレクトされる
				w.Write([]byte(LoadFixture("personal/facility/index.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		err := mockServer.Client.Login(&LoginParams{
			UserID:   "test_user",
			Password: "test_password",
		})

		if err != nil {
			t.Errorf("Expected successful login, got error: %v", err)
		}

		// リクエストの検証
		if len(mockServer.Requests) != 2 {
			t.Errorf("Expected 2 requests, got %d", len(mockServer.Requests))
		}
	})

	t.Run("AuthenticationFailure", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /index.aspx": func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(LoadFixture("index.html")))
			},
			"POST /index.aspx": func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(LoadFixture("index_failure_with_wrong_id_or_password.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		err := mockServer.Client.Login(&LoginParams{
			UserID:   "invalid_user",
			Password: "invalid_password",
		})

		if err != ErrAuthenticationFailed {
			t.Errorf("Expected authentication failure, got: %v", err)
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /index.aspx": func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(LoadFixture("index.html")))
			},
			"POST /index.aspx": func(w http.ResponseWriter, r *http.Request) {
				// サーバーエラーページを返す
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(LoadFixture("errorpage.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		err := mockServer.Client.Login(&LoginParams{
			UserID:   "test_user",
			Password: "test_password",
		})

		if err != ErrInternalServer {
			t.Errorf("Expected internal server error, got: %v", err)
		}
	})
}
