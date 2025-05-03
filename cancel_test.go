package tcmrsv

import (
	"net/http"
	"testing"
)

func TestCancelReservation(t *testing.T) {
	t.Run("SuccessfulCancellation", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /personal/facility/cancel.aspx": func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(LoadFixture("personal/facility/cancel.html")))
			},
			"POST /personal/facility/cancel.aspx": func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(LoadFixture("personal/facility/cancel_done.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		err := mockServer.Client.CancelReservation(&CancelReservationParams{
			ReservationID: "59253854-3628-f011-8c4e-000d3a51476f",
			Comment:       "テストのためキャンセル",
		})

		if err != nil {
			t.Errorf("Expected successful cancellation, got error: %v", err)
		}
	})

	t.Run("ValidationErrors", func(t *testing.T) {
		mockServer := NewMockServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// このハンドラーは入力検証エラーのため呼び出されるべきではない
			t.Errorf("Handler called despite validation errors")
		}))
		defer mockServer.Close()

		// 無効な予約IDのテスト
		err := mockServer.Client.CancelReservation(&CancelReservationParams{
			ReservationID: "invalid-id",
			Comment:       "テストのためキャンセル",
		})

		if err != ErrInvalidIDFormat {
			t.Errorf("Expected invalid ID format error, got: %v", err)
		}

		// 無効なコメントのテスト
		err = mockServer.Client.CancelReservation(&CancelReservationParams{
			ReservationID: "59253854-3628-f011-8c4e-000d3a51476f",
			Comment:       "",
		})

		if err != ErrInvalidComment {
			t.Errorf("Expected invalid comment error, got: %v", err)
		}
	})

	t.Run("CancellationFailure", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /personal/facility/cancel.aspx": func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(LoadFixture("personal/facility/cancel.html")))
			},
			"POST /personal/facility/cancel.aspx": func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(LoadFixture("personal/facility/cancel_failure_without_a_comment.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		err := mockServer.Client.CancelReservation(&CancelReservationParams{
			ReservationID: "59253854-3628-f011-8c4e-000d3a51476f",
			Comment:       "テストのためキャンセル",
		})

		if err != ErrCancelReservationFailed {
			t.Errorf("Expected cancellation failure, got: %v", err)
		}
	})

	t.Run("AuthenticationFailure", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /personal/facility/cancel.aspx": func(w http.ResponseWriter, r *http.Request) {
				// ログインページにリダイレクト
				w.Write([]byte(LoadFixture("index.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		err := mockServer.Client.CancelReservation(&CancelReservationParams{
			ReservationID: "59253854-3628-f011-8c4e-000d3a51476f",
			Comment:       "テストのためキャンセル",
		})

		if err != ErrAuthenticationFailed {
			t.Errorf("Expected authentication failure, got: %v", err)
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /personal/facility/cancel.aspx": func(w http.ResponseWriter, r *http.Request) {
				// サーバーエラーページを返す
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(LoadFixture("errorpage.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		err := mockServer.Client.CancelReservation(&CancelReservationParams{
			ReservationID: "59253854-3628-f011-8c4e-000d3a51476f",
			Comment:       "テストのためキャンセル",
		})

		if err != ErrInternalServer {
			t.Errorf("Expected internal server error, got: %v", err)
		}
	})
}
