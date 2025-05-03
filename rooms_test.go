package tcmrsv

import (
	"strings"
	"testing"
)

func TestGetRooms(t *testing.T) {
	client := New()
	rooms := client.GetRooms()

	// 部屋の数が正しいことを確認
	if len(rooms) == 0 {
		t.Error("Expected rooms to be returned, got empty slice")
	}

	// いくつかのサンプル部屋が存在することを確認
	roomIDs := make(map[string]bool)
	for _, room := range rooms {
		roomIDs[room.ID] = true
	}

	expectedRoomIDs := []string{
		"9ed14a61-a3e2-ef11-be20-7c1e52246dd3", // 楽屋2（U）旧楽屋201
		"29f2e624-2f48-ec11-8c60-002248696fd6", // 楽屋3（G）旧楽屋202
		"69f2e624-2f48-ec11-8c60-002248696fd6", // P 426（G）
	}

	for _, id := range expectedRoomIDs {
		if !roomIDs[id] {
			t.Errorf("Expected room with ID %s to exist", id)
		}
	}
}

func TestGetRoomsFiltered(t *testing.T) {
	client := New()

	t.Run("FilterByName", func(t *testing.T) {
		name := "楽屋"
		result := client.GetRoomsFiltered(GetRoomsFilteredParams{
			Name: &name,
		})

		if len(result) == 0 {
			t.Error("Expected rooms with name containing '楽屋', got none")
		}

		for _, room := range result {
			if !strings.Contains(room.Name, "楽屋") {
				t.Errorf("Room name '%s' does not contain '楽屋'", room.Name)
			}
		}
	})

	t.Run("FilterByID", func(t *testing.T) {
		id := "9ed14a61-a3e2-ef11-be20-7c1e52246dd3"
		result := client.GetRoomsFiltered(GetRoomsFilteredParams{
			ID: &id,
		})

		if len(result) != 1 {
			t.Errorf("Expected 1 room with ID %s, got %d", id, len(result))
		}

		if len(result) > 0 && result[0].ID != id {
			t.Errorf("Expected room ID %s, got %s", id, result[0].ID)
		}
	})

	t.Run("FilterByPianoType", func(t *testing.T) {
		result := client.GetRoomsFiltered(GetRoomsFilteredParams{
			PianoTypes: []RoomPianoType{RoomPianoTypeGrand},
		})

		if len(result) == 0 {
			t.Error("Expected rooms with grand pianos, got none")
		}

		for _, room := range result {
			if room.PianoType != RoomPianoTypeGrand {
				t.Errorf("Expected room with grand piano, got %s", room.PianoType)
			}
		}
	})

	t.Run("FilterByPianoNumber", func(t *testing.T) {
		result := client.GetRoomsFiltered(GetRoomsFilteredParams{
			PianoNumbers: []int{1},
		})

		if len(result) == 0 {
			t.Error("Expected rooms with piano number 1, got none")
		}

		for _, room := range result {
			if room.PianoNumber != 1 {
				t.Errorf("Expected room with piano number 1, got %d", room.PianoNumber)
			}
		}
	})

	t.Run("FilterByFloor", func(t *testing.T) {
		result := client.GetRoomsFiltered(GetRoomsFilteredParams{
			Floors: []int{4},
		})

		if len(result) == 0 {
			t.Error("Expected rooms on floor 4, got none")
		}

		for _, room := range result {
			if room.Floor != 4 {
				t.Errorf("Expected room on floor 4, got %d", room.Floor)
			}
		}
	})

	t.Run("FilterByIsBasement", func(t *testing.T) {
		isBasement := false
		result := client.GetRoomsFiltered(GetRoomsFilteredParams{
			IsBasement: &isBasement,
		})

		if len(result) == 0 {
			t.Error("Expected non-basement rooms, got none")
		}

		for _, room := range result {
			if room.IsBasement != isBasement {
				t.Errorf("Expected non-basement room, got basement room")
			}
		}
	})

	t.Run("FilterByCampus", func(t *testing.T) {
		result := client.GetRoomsFiltered(GetRoomsFilteredParams{
			Campuses: []Campus{CampusNakameguro},
		})

		if len(result) == 0 {
			t.Error("Expected rooms in Nakameguro campus, got none")
		}

		for _, room := range result {
			if room.Campus != CampusNakameguro {
				t.Errorf("Expected room in Nakameguro campus, got %s", room.Campus)
			}
		}
	})

	t.Run("MultipleFilters", func(t *testing.T) {
		result := client.GetRoomsFiltered(GetRoomsFilteredParams{
			PianoTypes: []RoomPianoType{RoomPianoTypeGrand},
			Floors:     []int{4},
			Campuses:   []Campus{CampusNakameguro},
		})

		if len(result) == 0 {
			t.Error("Expected rooms matching multiple filters, got none")
		}

		for _, room := range result {
			if room.PianoType != RoomPianoTypeGrand {
				t.Errorf("Expected room with grand piano, got %s", room.PianoType)
			}
			if room.Floor != 4 {
				t.Errorf("Expected room on floor 4, got %d", room.Floor)
			}
			if room.Campus != CampusNakameguro {
				t.Errorf("Expected room in Nakameguro campus, got %s", room.Campus)
			}
		}
	})

	t.Run("NoMatchingRooms", func(t *testing.T) {
		nonExistentID := "non-existent-id"
		result := client.GetRoomsFiltered(GetRoomsFilteredParams{
			ID: &nonExistentID,
		})

		if len(result) != 0 {
			t.Errorf("Expected no rooms with ID %s, got %d", nonExistentID, len(result))
		}
	})
}
