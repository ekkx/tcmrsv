package main

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/ekkx/tcmrsv"
)

type Config struct {
	UserID       string `env:"USER_ID"`
	UserPassword string `env:"USER_PW"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func main() {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}

	client := tcmrsv.New()

	if err := client.Login(&tcmrsv.LoginParams{
		UserID:   cfg.UserID,
		Password: cfg.UserPassword,
	}); err != nil {
		panic(err)
	}

	fmt.Println("Login successful!")

	if err = client.Reserve(&tcmrsv.ReserveParams{
		Campus:     tcmrsv.CampusNakameguro,
		RoomID:     "2df2e624-2f48-ec11-8c60-002248696fd6",
		Date:       tcmrsv.Today().AddDays(2),
		FromHour:   12,
		FromMinute: 0,
		ToHour:     14,
		ToMinute:   0,
	}); err != nil {
		panic(err)
	}

	fmt.Println("Reservation successful!")

	reservations, err := client.GetMyReservations()
	if err != nil {
		panic(err)
	}

	for _, reservation := range reservations {
		fmt.Println(reservation)
	}

	fmt.Println("Total reservations:", len(reservations))
}
