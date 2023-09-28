package room

import (
	"errors"
	"log"

	"github.com/size12/planning-poker/internal/entity"
	"github.com/size12/planning-poker/internal/entity/update"
)

func (r *Room) Update(u *update.Update) error {
	switch u.Type {
	case update.Connected:
		ID := u.PlayerID
		player, err := entity.NewPlayer(ID)
		if err != nil {
			log.Printf("Failed create new player with id %s: %v\n", ID, err)
			return err
		}

		err = r.AddPlayer(player)
		if err != nil {
			log.Printf("Failed add player %v to room %s: %v\n", player, r.ID, err)
			return err
		}
	case update.Vote:
		ID := u.PlayerID

		err := r.Vote(ID, u.Vote)

		if errors.Is(err, ErrPlayerNotFound) {
			log.Printf("Got vote from non-existent user with id: %s in room %s\n", ID, r.ID)
			return ErrPlayerNotFound
		}

		if err != nil {
			log.Printf("Failed set player %s vote %s in room %s: %v\n", ID, *u.Vote, r.ID, err)
			return err
		}

		if r.CountVotes() == r.Size() {
			err = r.NextStatus(r.adminID)
			if err != nil {
				return err
			}
			log.Printf("All players in room %s voted, revealing\n", r.ID)
		}

	case update.ChangeStatus:
		err := r.NextStatus(u.PlayerID)
		if err != nil {
			log.Printf("Failed start new voting in room %s: %v\n", r.ID, err)
			return err
		}

	case update.ChangePlayerStatus:
		err := r.NextPlayerStatus(u.PlayerID)
		if err != nil {
			log.Printf("Failed start change voting type in room %s: %v\n", r.ID, err)
			return err
		}

		if r.CountVotes() == r.Size() {
			err = r.NextStatus(r.adminID)
			if err != nil {
				return err
			}
			log.Printf("All players in room %s voted, revealing\n", r.ID)
		}
	case update.ChangePlayerName:
		err := r.ChangePlayerName(u.PlayerID, u.Message)
		if err != nil {
			log.Printf("Failed change player %v name in room %v: %v\n", r.ID, u.PlayerID, err)
			return err
		}

		return nil
	}

	return nil
}
