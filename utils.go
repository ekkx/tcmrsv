package tcmrsv

import "time"

func jst() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}
