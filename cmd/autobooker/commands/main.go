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

	// err = rsv.Reserve(&tcmrsv.ReserveParams{
	// 	Campus:     tcmrsv.CampusIkebukuro,
	// 	RoomID:     "bbf2e624-2f48-ec11-8c60-002248696fd6",
	// 	Date:       time.Now().AddDate(0, 0, 2),
	// 	FromHour:   17,
	// 	FromMinute: 00,
	// 	ToHour:     22,
	// 	ToMinute:   30,
	// })
	// if err != nil {
	// 	return err
	// }

	// err = rsv.Reserve(&tcmrsv.ReserveParams{
	// 	Campus:     tcmrsv.CampusNakameguro,
	// 	RoomID:     "81f2e624-2f48-ec11-8c60-002248696fd6",
	// 	Date:       time.Now().AddDate(0, 0, 1),
	// 	FromHour:   14,
	// 	FromMinute: 00,
	// 	ToHour:     14,
	// 	ToMinute:   30,
	// })
	// if err != nil {
	// 	return err
	// }

	rooms, err := rsv.GetRoomAvailability(&tcmrsv.GetRoomAvailabilityParams{
		Campus: tcmrsv.CampusNakameguro,
		Date:   time.Now().In(tcmrsv.JST()),
	})
	if err != nil {
		return err
	}

	for _, room := range rooms {
		if room.Room.Name == "C401" {
			fmt.Println("--------")
			fmt.Println(room.Room.Name)
			fmt.Println(room.Room.PianoType)
			fmt.Println(room.AvailableTimes)
		}
	}

	// rsvs, err := rsv.GetMyReservations()
	// if err != nil {
	// 	return err
	// }

	// for _, r := range rsvs {
	// 	fmt.Println("--------")
	// 	fmt.Println(r)
	// }

	return nil
}
