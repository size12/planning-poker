package room

import (
	"github.com/google/uuid"
	"github.com/size12/planning-poker/internal/entity"
)

func (r *Room) AddPlayer(player *entity.Player) error {
	r.Lock()
	defer r.Unlock()

	r.Players = append(r.Players, player)
	return nil
}

func (r *Room) RemovePlayer(ID uuid.UUID) error {
	r.Lock()
	defer r.Unlock()

	for index, player := range r.Players {
		if player.ID == ID {
			r.Players = append(r.Players[:index], r.Players[index+1:]...)
			return nil
		}
	}
	return ErrPlayerNotFound
}

func (r *Room) PlayerByID(ID uuid.UUID) (*entity.Player, error) {
	r.RLock()
	defer r.RUnlock()

	for _, player := range r.Players {
		if player.ID == ID {
			return player, nil
		}
	}

	return nil, ErrPlayerNotFound
}

func (r *Room) ChangePlayerName(ID uuid.UUID, name string) error {
	r.RLock()
	defer r.RUnlock()

	for _, player := range r.Players {
		if player.ID == ID {
			return player.SetName(name)
		}
	}

	return ErrPlayerNotFound
}
