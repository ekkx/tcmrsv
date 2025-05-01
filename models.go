package tcmrsv

type Campus string

const (
	CampusUnknown    Campus = "-1"
	CampusIkebukuro  Campus = "1"
	CampusNakameguro Campus = "2"
)

func (c Campus) IsValid() bool {
	switch c {
	case CampusIkebukuro, CampusNakameguro:
		return true
	default:
		return false
	}
}

type RoomPianoType string

const (
	RoomPianoTypeGrand   RoomPianoType = "grand"
	RoomPianoTypeUpright RoomPianoType = "upright"
	RoomPianoTypeUnknown RoomPianoType = "unknown"
	RoomPianoTypeNone    RoomPianoType = "none"
)

func (r RoomPianoType) IsValid() bool {
	switch r {
	case RoomPianoTypeGrand, RoomPianoTypeUpright, RoomPianoTypeUnknown, RoomPianoTypeNone:
		return true
	default:
		return false
	}
}

type Room struct {
	ID          string
	Name        string
	PianoType   RoomPianoType
	PianoNumber int
	IsClassroom bool
	IsBasement  bool
	Campus      Campus
	Floor       int
}

type AvailableTime struct {
	Hour   int
	Minute int
}

type RoomAvailability struct {
	Room           Room
	AvailableTimes []AvailableTime
}

type Reservation struct {
	ID         string
	Campus     Campus
	CampusName string
	Date       string
	RoomName   string
	TimeRange  string
}
