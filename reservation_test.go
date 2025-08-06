package tcmrsv

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestGetMyReservations(t *testing.T) {
	t.Run("SuccessfulRetrieval", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /personal/facility/index.aspx": func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(LoadFixture("personal/facility/index.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		reservations, err := mockServer.Client.GetMyReservations()

		if err != nil {
			t.Errorf("Expected successful retrieval, got error: %v", err)
		}

		// 予約の検証
		if len(reservations) != 2 {
			t.Errorf("Expected 2 reservations, got %d", len(reservations))
		}

		// 予約内容の検証
		if reservations[0].Campus != CampusIkebukuro {
			t.Errorf("Expected campus to be Ikebukuro, got %v", reservations[0].Campus)
		}

		if reservations[0].Date != NewDate(2025, 5, 5) {
			t.Errorf("Expected date to be 2025年05月05日（月）, got %s", reservations[0].Date)
		}

		if reservations[0].FromHour == 17 && reservations[0].FromMinute == 0 && reservations[0].ToHour == 22 && reservations[0].ToMinute == 30 {
			t.Errorf("Expected time range to be 17:00-22:30, got %s", fmt.Sprintf("%02d:%02d-%02d:%02d", reservations[0].FromHour, reservations[0].FromMinute, reservations[0].ToHour, reservations[0].ToMinute))
		}

		if reservations[0].RoomName != "A414（G）" {
			t.Errorf("Expected room name to be A414（G）, got %s", reservations[0].RoomName)
		}

		if reservations[0].ID != "fa791156-cc27-f011-8c4e-000d3ace9c3e" {
			t.Errorf("Expected ID to be fa791156-cc27-f011-8c4e-000d3ace9c3e, got %s", reservations[0].ID)
		}
	})

	t.Run("EmptyReservations", func(t *testing.T) {
		// 予約がない場合のHTMLを作成
		originalHTML := LoadFixture("personal/facility/index.html")

		// 予約リストの部分を空にする
		startMarker := `<div id="reservation-list">`
		endMarker := `</div>`

		startIndex := strings.Index(originalHTML, startMarker)
		if startIndex == -1 {
			t.Fatalf("Could not find start marker in HTML")
		}

		// 最初の予約リストの終了タグを見つける
		endIndex := strings.Index(originalHTML[startIndex:], endMarker)
		if endIndex == -1 {
			t.Fatalf("Could not find end marker in HTML")
		}
		endIndex += startIndex + len(endMarker)

		// 予約リストを空にする
		emptyHTML := originalHTML[:startIndex] + startMarker + endMarker + originalHTML[endIndex:]

		routes := map[string]http.HandlerFunc{
			"GET /personal/facility/index.aspx": func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(emptyHTML))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		reservations, err := mockServer.Client.GetMyReservations()

		if err != nil {
			t.Errorf("Expected successful retrieval, got error: %v", err)
		}

		if len(reservations) != 0 {
			t.Errorf("Expected 0 reservations, got %d", len(reservations))
		}
	})

	t.Run("AuthenticationFailure", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /personal/facility/index.aspx": func(w http.ResponseWriter, r *http.Request) {
				// ログインページにリダイレクト
				w.Write([]byte(LoadFixture("index.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		_, err := mockServer.Client.GetMyReservations()

		if err != ErrAuthenticationFailed {
			t.Errorf("Expected authentication failure, got: %v", err)
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /personal/facility/index.aspx": func(w http.ResponseWriter, r *http.Request) {
				// サーバーエラーページを返す
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(LoadFixture("errorpage.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		_, err := mockServer.Client.GetMyReservations()

		if err != ErrInternalServer {
			t.Errorf("Expected internal server error, got: %v", err)
		}
	})
}

func TestReserve(t *testing.T) {
	t.Run("SuccessfulReservation", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /personal/facility/confirms.aspx": func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(LoadFixture("personal/facility/confirms.html")))
			},
			"POST /personal/facility/confirms.aspx": func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(LoadFixture("personal/facility/done.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		tomorrow := Today().AddDays(1)

		err := mockServer.Client.Reserve(&ReserveParams{
			Campus:     CampusNakameguro,
			RoomID:     "23f2e624-2f48-ec11-8c60-002248696fd6", // P 200（G）
			Date:       tomorrow,
			FromHour:   10,
			FromMinute: 30,
			ToHour:     11,
			ToMinute:   30,
		})

		if err != nil {
			t.Errorf("Expected successful reservation, got error: %v", err)
		}
	})

	t.Run("ValidationErrors", func(t *testing.T) {
		mockServer := NewMockServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// このハンドラーは入力検証エラーのため呼び出されるべきではない
			t.Errorf("Handler called despite validation errors")
		}))
		defer mockServer.Close()

		// 無効なキャンパスのテスト
		err := mockServer.Client.Reserve(&ReserveParams{
			Campus:     CampusUnknown,
			RoomID:     "23f2e624-2f48-ec11-8c60-002248696fd6",
			Date:       Today().AddDays(1),
			FromHour:   10,
			FromMinute: 30,
			ToHour:     11,
			ToMinute:   30,
		})

		if err != ErrInvalidCampus {
			t.Errorf("Expected invalid campus error, got: %v", err)
		}

		// 無効な部屋IDのテスト
		err = mockServer.Client.Reserve(&ReserveParams{
			Campus:     CampusNakameguro,
			RoomID:     "invalid-id",
			Date:       Today().AddDays(1),
			FromHour:   10,
			FromMinute: 30,
			ToHour:     11,
			ToMinute:   30,
		})

		if err != ErrInvalidIDFormat {
			t.Errorf("Expected invalid ID format error, got: %v", err)
		}

		// 範囲外の日付のテスト
		futureDate := Today().AddDays(4)
		err = mockServer.Client.Reserve(&ReserveParams{
			Campus:     CampusNakameguro,
			RoomID:     "23f2e624-2f48-ec11-8c60-002248696fd6",
			Date:       futureDate,
			FromHour:   10,
			FromMinute: 30,
			ToHour:     11,
			ToMinute:   30,
		})

		if err != ErrDateOutOfRange {
			t.Errorf("Expected date out of range error, got: %v", err)
		}

		// 無効な時間範囲のテスト
		err = mockServer.Client.Reserve(&ReserveParams{
			Campus:     CampusNakameguro,
			RoomID:     "23f2e624-2f48-ec11-8c60-002248696fd6",
			Date:       Today().AddDays(1),
			FromHour:   11,
			FromMinute: 30,
			ToHour:     10, // 終了時間が開始時間より前
			ToMinute:   30,
		})

		if err != ErrInvalidTimeRange {
			t.Errorf("Expected invalid time range error, got: %v", err)
		}
	})

	t.Run("ReservationFailure", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /personal/facility/confirms.aspx": func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(LoadFixture("personal/facility/confirms.html")))
			},
			"POST /personal/facility/confirms.aspx": func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(LoadFixture("personal/facility/confirms_failure.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		tomorrow := Today().AddDays(1)

		err := mockServer.Client.Reserve(&ReserveParams{
			Campus:     CampusNakameguro,
			RoomID:     "23f2e624-2f48-ec11-8c60-002248696fd6",
			Date:       tomorrow,
			FromHour:   10,
			FromMinute: 30,
			ToHour:     11,
			ToMinute:   30,
		})

		if err != ErrCreateReservationFailed {
			t.Errorf("Expected reservation failure, got: %v", err)
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		routes := map[string]http.HandlerFunc{
			"GET /personal/facility/confirms.aspx": func(w http.ResponseWriter, r *http.Request) {
				// サーバーエラーページを返す
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(LoadFixture("errorpage.html")))
			},
		}

		mockServer := NewMockServer(CreateHandler(routes))
		defer mockServer.Close()

		tomorrow := Today().AddDays(1)

		err := mockServer.Client.Reserve(&ReserveParams{
			Campus:     CampusNakameguro,
			RoomID:     "23f2e624-2f48-ec11-8c60-002248696fd6",
			Date:       tomorrow,
			FromHour:   10,
			FromMinute: 30,
			ToHour:     11,
			ToMinute:   30,
		})

		if err != ErrInternalServer {
			t.Errorf("Expected internal server error, got: %v", err)
		}
	})
}
