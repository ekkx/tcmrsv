package tcmrsv

import (
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type GetRoomAvailabilityParams struct {
	Campus Campus
	Date   time.Time
}

func (rsv *TCMRSV) GetRoomAvailability(params *GetRoomAvailabilityParams) ([]Room, error) {
	u, err := url.Parse(ENDPOINT_RESERVE)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("campus", string(params.Campus))
	q.Set("ymd", params.Date.Format("2006/01/02 15:04:05"))
	u.RawQuery = q.Encode()

	res, err := rsv.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	z := html.NewTokenizer(res.Body)

	var rooms []Room
	var insideTable bool
	var insideRow bool
	var colIndex int
	var currentRoom *Room

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return rooms, nil
			}
			return nil, z.Err()

		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()

			switch t.Data {
			case "table":
				for _, attr := range t.Attr {
					if attr.Key == "id" && strings.HasPrefix(attr.Val, "aspTable") {
						insideTable = true
					}
				}

			case "tr":
				if insideTable {
					colIndex = 0
					insideRow = true
				}

			case "td":
				if insideRow {
					colIndex++
					attrs := make(map[string]string)
					for _, a := range t.Attr {
						attrs[a.Key] = a.Val
					}

					// 部屋
					if colIndex == 1 {
						class := attrs["class"]
						if class == "g" || class == "a" {
							currentRoom = &Room{
								Type: RoomType(class),
							}
						}
					}

					// 空き状況
					if colIndex >= 2 && currentRoom != nil {
						class := attrs["class"]
						if strings.HasPrefix(class, "judgment") || attrs["style"] == "text-align:center;" {
							var roomID string
							var hour int
							var minute int
							hasInput := false

							for {
								tt2 := z.Next()
								if tt2 == html.StartTagToken || tt2 == html.SelfClosingTagToken {
									t2 := z.Token()
									if t2.Data == "input" {
										inputAttrs := make(map[string]string)
										for _, a := range t2.Attr {
											inputAttrs[a.Key] = a.Val
										}
										id, ok := inputAttrs["id"]
										if ok {
											parts := strings.Split(id, ",")
											if len(parts) == 2 {
												roomID = parts[0]
												timeStr := parts[1]
												if len(timeStr) == 4 {
													hour, _ = strconv.Atoi(timeStr[:2])
													minute, _ = strconv.Atoi(timeStr[2:])
												}
												hasInput = true
											}
										}
										break
									}
								} else if tt2 == html.EndTagToken && z.Token().Data == "td" {
									break
								}
							}

							if hasInput {
								if currentRoom.ID == "" {
									currentRoom.ID = roomID
								}
								currentRoom.Slots = append(currentRoom.Slots, TimeSlot{
									Hour:       hour,
									Minute:     minute,
									Reservable: class == "judgment4",
								})
							}
						}
					}
				}

			case "span":
				if insideRow && colIndex == 1 && currentRoom != nil {
					tt2 := z.Next()
					if tt2 == html.TextToken {
						currentRoom.Name = strings.TrimSpace(z.Token().Data)
					}
				}
			}

		case html.EndTagToken:
			t := z.Token()
			switch t.Data {
			case "tr":
				if insideRow {
					if currentRoom != nil && currentRoom.Name != "" {
						rooms = append(rooms, *currentRoom)
						currentRoom = nil
					}
					insideRow = false
				}
			case "table":
				insideTable = false
			}
		}
	}
}
