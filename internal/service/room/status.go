package room

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/size12/planning-poker/internal/entity/voting"
)

func (r *Room) setStatus(status voting.Status) error {
	r.Lock()
	defer r.Unlock()

	r.Status = status
	return nil
}

func (r *Room) NextStatus(ID uuid.UUID) error {
	if !r.IsAdmin(ID) {
		return ErrAccessDenied
	}

	var err error

	switch r.Status {
	case voting.RoomStatusWaiting:
		err = r.setStatus(voting.RoomStatusVoting)
	case voting.RoomStatusVoting:
		err = r.setStatus(voting.RoomStatusRevealed)
	case voting.RoomStatusRevealed:
		err = r.setStatus(voting.RoomStatusVoting)
	default:
		return fmt.Errorf("unknown status %s", r.Status)
	}

	if err != nil {
		return err
	}

	if r.Status == voting.RoomStatusVoting {
		if err := r.ClearVotes(ID); err != nil {
			return err
		}
	}

	return nil
}

func (r *Room) NextPlayerStatus(ID uuid.UUID) error {
	player, err := r.PlayerByID(ID)

	if err != nil {
		return err
	}

	switch player.Status {
	case voting.PlayerVoting:
		player.SetVotingStatus(voting.PlayerObserving)
	case voting.PlayerObserving:
		player.SetVotingStatus(voting.PlayerVoting)
	default:
		return fmt.Errorf("unknown voting status %v", player.Status)
	}

	return nil
}
