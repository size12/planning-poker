package room

import (
	"github.com/google/uuid"
	"github.com/size12/planning-poker/internal/entity/voting"
)

func (r *Room) Vote(playerID uuid.UUID, vote *voting.Vote) error {
	r.Lock()
	defer r.Unlock()

	for _, player := range r.Players {
		if player.ID == playerID {
			player.Vote(vote)
			return nil
		}
	}

	return ErrPlayerNotFound
}

// ClearVotes set all votes to nil. Only admin can do this.
func (r *Room) ClearVotes(ID uuid.UUID) error {
	if !r.IsAdmin(ID) {
		return ErrAccessDenied
	}

	for _, player := range r.Players {
		err := player.ClearVote()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Room) CountVotes() int {
	r.RLock()
	defer r.RUnlock()

	cnt := 0

	for _, player := range r.Players {
		if player.Voted != nil || player.Status == voting.PlayerObserving {
			cnt += 1
		}
	}

	return cnt
}
