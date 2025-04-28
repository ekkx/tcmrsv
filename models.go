package tcmrsv

type Campus string

const (
	CampusUnknown    Campus = "-1"
	CampusIkebukuro  Campus = "1"
	CampusNakameguro Campus = "2"
)

type RoomType string

const (
	RoomTypeGrand   RoomType = "g"
	RoomTypeUpright RoomType = "a"
)

type TimeSlot struct {
	Hour       int
	Minute     int
	Reservable bool
}

type Room struct {
	ID    string
	Name  string
	Type  RoomType
	Slots []TimeSlot
}

type Reservation struct {
	ID         string
	Campus     Campus
	CampusName string
	Date       string
	RoomName   string
	TimeRange  string
}
