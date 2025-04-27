# tcmrsv

練習室自動予約ツール

### Getting Started

```go
package main

import (
	"time"

	"github.com/ekkx/tcmrsv"
)

func main() {
	rsv := tcmrsv.New()

	if _, err := rsv.Login(&tcmrsv.LoginParams{
		ID:       "-----",
		Password: "-----",
	}); err != nil {
		panic(err)
	}

	if _, err := rsv.Reserve(&tcmrsv.ReserveParams{
		Campus:     tcmrsv.CampusIkebukuro,
		RoomID:     "42d1eacc-60d5-428b-8c64-aef11a512c30",
		Start:      time.Date(2025, 4, 29, 0, 0, 0, 0, time.Local),
		FromHour:   12,
		FromMinute: 0,
		ToHour:     14,
		ToMinute:   0,
	}); err != nil {
		panic(err)
	}
}
```
