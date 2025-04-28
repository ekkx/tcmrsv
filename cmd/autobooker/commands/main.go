package commands

import (
	"fmt"

	"github.com/ekkx/tcmrsv"
)

func Run() error {
	cfg, err := NewConfig()
	if err != nil {
		return err
	}

	rsv := tcmrsv.New()

	_, err = rsv.Login(&tcmrsv.LoginParams{
		UserID:   cfg.UserID,
		Password: cfg.UserPassword,
	})
	if err != nil {
		return err
	}

	// rooms, err := rsv.GetRoomAvailability(&tcmrsv.GetRoomAvailabilityParams{
	// 	Campus: tcmrsv.CampusIkebukuro,
	// 	Date:   time.Now().AddDate(0, 0, 1).Local(),
	// })
	// if err != nil {
	// 	return err
	// }

	// for _, room := range rooms {
	// 	fmt.Println("--------")
	// 	fmt.Println(room.Name)
	// 	fmt.Println(room.Type == tcmrsv.RoomTypeGrand)
	// }

	// if _, err := rsv.Reserve(&tcmrsv.ReserveParams{
	// 	Campus:     tcmrsv.CampusIkebukuro,
	// 	RoomID:     "42d1eacc-60d5-428b-8c64-aef11a512c30",
	// 	Start:      time.Date(2025, 4, 29, 0, 0, 0, 0, time.Local),
	// 	FromHour:   12,
	// 	FromMinute: 0,
	// 	ToHour:     13,
	// 	ToMinute:   0,
	// }); err != nil {
	// 	panic(err)
	// }

	rsvs, err := rsv.GetMyReservations()
	if err != nil {
		return err
	}

	for _, r := range rsvs {
		fmt.Println("--------")
		fmt.Println(r)
	}

	var ikbkrRsvs []tcmrsv.Reservation
	for _, r := range rsvs {
		if r.Campus == tcmrsv.CampusIkebukuro {
			ikbkrRsvs = append(ikbkrRsvs, r)
		}
	}

	_, err = rsv.CancelReservation(&tcmrsv.CancelReservationParams{
		ReservationID: ikbkrRsvs[0].ID,
		Comment:       "間違えました。",
	})
	if err != nil {
		return err
	}

	rsvs, err = rsv.GetMyReservations()
	if err != nil {
		return err
	}

	for _, r := range rsvs {
		fmt.Println("--------")
		fmt.Println(r)
	}

	return nil
}
