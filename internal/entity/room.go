package entity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/size12/planning-poker/internal/entity/voting"
)

var ErrPlayerNotFound = errors.New("player with such ID not found")

type Room struct {
	id      uuid.UUID
	name    string
	players []*Player
	status  voting.Status
}

func NewRoom(name string) (*Room, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	players := make([]*Player, 0, 5)
	return &Room{
		id:      id,
		name:    name,
		players: players,
	}, nil
}

func (r *Room) ID() (uuid.UUID, error) {
	return r.id, nil
}

func (r *Room) Name() (string, error) {
	return r.name, nil
}

func (r *Room) AddPlayer(player *Player) error {
	r.players = append(r.players, player)
	return nil
}

func (r *Room) RemovePlayer(ID uuid.UUID) error {
	for index, player := range r.players {
		if player.ID() == ID {
			r.players = append(r.players[:index], r.players[index+1:]...)
			return nil
		}
	}
	return ErrPlayerNotFound
}

func (r *Room) RemoveAllPlayers() error {
	r.players = make([]*Player, 0)
	return nil
}

func (r *Room) ClearVotes() error {
	for _, player := range r.players {
		err := player.ClearVote()
		if err != nil {
			return err
		}
	}
	return nil
}
