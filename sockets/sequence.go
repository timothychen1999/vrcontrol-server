package sockets

import "sort"

type SequenceUpdate struct {
	Player   *Player
	Sequence int
}

func (r *Room) PlayerSequenceUpdate() []SequenceUpdate {
	names := make([]string, 0, len(r.Players))
	maps := make(map[string]*Player, len(r.Players))
	used_sequences := make(map[int]bool, len(r.Players))
	for player := range r.Players {
		if player == nil {
			continue
		}
		if seq, ok := r.AssignedSequence[player.DeiviceID]; ok {
			used_sequences[seq] = true
			// Player already has a sequence assigned, skip
			continue
		}
		names = append(names, player.DeiviceID)
		maps[player.DeiviceID] = player
	}

	// Sort by lexicographic order
	if len(names) == 0 {
		return nil
	}
	sequenceUpdates := make([]SequenceUpdate, 0, len(r.Players))
	sort.Strings(names)
	i := 0
	for _, name := range names {
		if _, ok := used_sequences[i]; ok {
			// Sequence already used, skip
			i++
			continue
		}
		player := maps[name]
		if player == nil {
			continue
		}
		// Player sequence use 0-based index
		sequenceUpdates = append(sequenceUpdates, SequenceUpdate{
			Player:   player,
			Sequence: i,
		})
		player.Sequence = i
		i++
	}
	return sequenceUpdates

}
