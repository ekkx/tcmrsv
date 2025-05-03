package tcmrsv

import "time"

func JST() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}
