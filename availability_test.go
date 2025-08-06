package tcmrsv

import (
	"net/http"
	"testing"
)

func TestGetRoomAvailability(t *testing.T) {
	t.Run("SuccessfulAvailabilityRetrievalWithInputs", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /personal/facility/reserve.aspx": func(w http.ResponseWriter, r *http.Request) {
				// クエリパラメータの検証
				q := r.URL.Query()
				campus := q.Get("campus")
				ymd := q.Get("ymd")

				if campus != string(CampusNakameguro) {
					t.Errorf("Expected campus to be %s, got %s", CampusNakameguro, campus)
				}

				if ymd == "" {
					t.Error("Expected ymd parameter to be set")
				}

				// 練習室の空き状況ページのHTMLを返す
				w.Write([]byte(LoadFixture("personal/facility/reserve_with_inputs.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		// 現在時刻から1日後の日付を設定
		tomorrow := Today().AddDays(1)

		availabilities, err := mockServer.Client.GetRoomAvailability(&GetRoomAvailabilityParams{
			Campus: CampusNakameguro,
			Date:   tomorrow,
		})

		if err != nil {
			t.Errorf("Expected successful availability retrieval, got error: %v", err)
		}

		if len(availabilities) == 0 {
			t.Error("Expected room availabilities to be returned, got empty slice")
		}

		// 少なくとも1つの部屋に空き時間があることを確認
		hasAvailableTimes := false
		for _, avail := range availabilities {
			if len(avail.AvailableTimes) > 0 {
				hasAvailableTimes = true
				break
			}
		}

		if !hasAvailableTimes {
			t.Error("Expected at least one room to have available times")
		}

		// リクエストの検証
		if len(mockServer.Requests) != 1 {
			t.Errorf("Expected 1 request, got %d", len(mockServer.Requests))
		}
	})

	t.Run("SuccessfulAvailabilityRetrievalWithoutInputs", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /personal/facility/reserve.aspx": func(w http.ResponseWriter, r *http.Request) {
				// クエリパラメータの検証
				q := r.URL.Query()
				campus := q.Get("campus")
				ymd := q.Get("ymd")

				if campus != string(CampusIkebukuro) {
					t.Errorf("Expected campus to be %s, got %s", CampusIkebukuro, campus)
				}

				if ymd == "" {
					t.Error("Expected ymd parameter to be set")
				}

				// 練習室の空き状況ページのHTMLを返す
				w.Write([]byte(LoadFixture("personal/facility/reserve_without_inputs.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		// 現在時刻から1日後の日付を設定
		tomorrow := Today().AddDays(1)

		availabilities, err := mockServer.Client.GetRoomAvailability(&GetRoomAvailabilityParams{
			Campus: CampusIkebukuro,
			Date:   tomorrow,
		})

		if err != nil {
			t.Errorf("Expected successful availability retrieval, got error: %v", err)
		}

		if len(availabilities) == 0 {
			t.Error("Expected room availabilities to be returned, got empty slice")
		}

		// 少なくとも1つの部屋に空き時間があることを確認
		hasAvailableTimes := false
		for _, avail := range availabilities {
			if len(avail.AvailableTimes) > 0 {
				hasAvailableTimes = true
				break
			}
		}

		if !hasAvailableTimes {
			t.Error("Expected at least one room to have available times")
		}

		// リクエストの検証
		if len(mockServer.Requests) != 1 {
			t.Errorf("Expected 1 request, got %d", len(mockServer.Requests))
		}
	})

	t.Run("InvalidCampus", func(t *testing.T) {
		client := New()

		// 現在時刻から1日後の日付を設定
		tomorrow := Today().AddDays(1)

		_, err := client.GetRoomAvailability(&GetRoomAvailabilityParams{
			Campus: CampusUnknown,
			Date:   tomorrow,
		})

		if err != ErrInvalidCampus {
			t.Errorf("Expected invalid campus error, got: %v", err)
		}
	})

	t.Run("InvalidDateRange", func(t *testing.T) {
		client := New()

		// 現在時刻から3日後の日付を設定（2日以内の制限を超える）
		futureDateOutOfRange := Today().AddDays(3)

		_, err := client.GetRoomAvailability(&GetRoomAvailabilityParams{
			Campus: CampusNakameguro,
			Date:   futureDateOutOfRange,
		})

		if err != ErrInvalidTimeRange {
			t.Errorf("Expected invalid time range error, got: %v", err)
		}
	})

	t.Run("AuthenticationFailure", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /personal/facility/reserve.aspx": func(w http.ResponseWriter, r *http.Request) {
				// 認証エラーの場合はログインページにリダイレクトされる
				w.Write([]byte(LoadFixture("index.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		// 現在時刻から1日後の日付を設定
		tomorrow := Today().AddDays(1)

		_, err := mockServer.Client.GetRoomAvailability(&GetRoomAvailabilityParams{
			Campus: CampusNakameguro,
			Date:   tomorrow,
		})

		if err != ErrAuthenticationFailed {
			t.Errorf("Expected authentication failure, got: %v", err)
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /personal/facility/reserve.aspx": func(w http.ResponseWriter, r *http.Request) {
				// サーバーエラーページを返す
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(LoadFixture("errorpage.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		// 現在時刻から1日後の日付を設定
		tomorrow := Today().AddDays(1)

		_, err := mockServer.Client.GetRoomAvailability(&GetRoomAvailabilityParams{
			Campus: CampusNakameguro,
			Date:   tomorrow,
		})

		if err != ErrInternalServer {
			t.Errorf("Expected internal server error, got: %v", err)
		}
	})
}
