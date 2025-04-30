package commands

import (
	"fmt"
	"time"

	"github.com/ekkx/tcmrsv"
)

func Run() error {
	cfg, err := NewConfig()
	if err != nil {
		return err
	}

	rsv := tcmrsv.New()

	err = rsv.Login(&tcmrsv.LoginParams{
		UserID:   cfg.UserID,
		Password: cfg.UserPassword,
	})
	if err != nil {
		return err
	}

	// err = rsv.CancelReservation(&tcmrsv.CancelReservationParams{
	// 	ReservationID: "8de87515-b625-f011-8c4e-000d3a51476f",
	// 	Comment:       "a",
	// })
	// if err != nil {
	// 	return err
	// }

	// err = rsv.Reserve(&tcmrsv.ReserveParams{
	// 	Campus:     tcmrsv.CampusIkebukuro,
	// 	RoomID:     "b9f2e624-2f48-ec11-8c60-002248696fd6",
	// 	Date:       time.Now().AddDate(0, 0, 1),
	// 	FromHour:   20,
	// 	FromMinute: 00,
	// 	ToHour:     22,
	// 	ToMinute:   30,
	// })
	// if err != nil {
	// 	return err
	// }

	rooms, err := rsv.GetRoomAvailability(&tcmrsv.GetRoomAvailabilityParams{
		Campus: tcmrsv.CampusIkebukuro,
		Date:   time.Now().AddDate(0, 0, 1),
	})
	if err != nil {
		return err
	}

	for _, room := range rooms {
		fmt.Println("--------")
		fmt.Println(room)
	}

	rsvs, err := rsv.GetMyReservations()
	if err != nil {
		return err
	}

	for _, r := range rsvs {
		fmt.Println("--------")
		fmt.Println(r)
	}

	return nil
}
