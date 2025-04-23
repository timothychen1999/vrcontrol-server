package logic

import "github.com/timothychen1999/vrcontrol-server/sockets"

var (
	ServerBlockChapter = []int{}
	NoSyncChapter      = []int{0, 1, 2, 3, 4, 5, 6, 7, 9}
)

func MovementCheck(r *sockets.Room, from int) (sockets.Movement, bool) {
	if r == nil {
		panic("room is nil")
	}
	if len(r.Players) == 0 {
		return sockets.Movement{}, false
	}
	// Check if sync is not required
	for _, chapter := range ServerBlockChapter {
		if chapter == from {
			return sockets.Movement{}, false
		}
	}
	for _, chapter := range NoSyncChapter {
		if chapter == from {
			return sockets.Movement{
				Force:            false,
				DestinationStage: from + 1,
				Target:           "all",
				Broadcast:        false,
			}, true
		}
	}
	// Player synconization chapter, require all players to be ready
	for player := range r.Players {
		if player == nil {
			continue
		}
		if player.Stage > from || (player.ReadyToMove && player.Stage == from) {
			continue
		}
		// Player is not ready to move, wait for them
		return sockets.Movement{}, false
	}
	return sockets.Movement{
		Force:            false,
		DestinationStage: from + 1,
		Target:           "all",
		Broadcast:        true,
	}, true
}
