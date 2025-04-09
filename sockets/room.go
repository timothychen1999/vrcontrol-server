package sockets

type Room struct {
	RoomID           string
	PlayerBroadcast  chan []byte
	PlayerRegister   chan *Player
	PlayerUnregister chan *Player
	Players          map[*Player]bool
}

func NewRoom(roomID string) *Room {
	room := &Room{
		RoomID:           roomID,
		PlayerBroadcast:  make(chan []byte, 1024),
		PlayerRegister:   make(chan *Player),
		PlayerUnregister: make(chan *Player),
		Players:          make(map[*Player]bool),
	}
	return room
}
func (r *Room) Run() {
	for {
		select {
		case player := <-r.PlayerRegister:
			r.Players[player] = true
		case player := <-r.PlayerUnregister:
			if _, ok := r.Players[player]; ok {
				delete(r.Players, player)
				close(player.InChannel)
			}
		case message := <-r.PlayerBroadcast:
			//Handle Messages from Players
			println("Message Received: ", string(message))
		}
	}
}
func ConnectToRoom() {

}
