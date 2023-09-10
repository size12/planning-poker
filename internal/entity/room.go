package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/size12/planning-poker/internal/entity/voting"
)

var ErrPlayerNotFound = errors.New("player with such ID does not exists")
var ErrEmptyRoomName = errors.New("room can't have empty name")

type Room struct {
	ID      uuid.UUID     `json:"id"`
	Name    string        `json:"name"`
	Players []*Player     `json:"players"`
	Status  voting.Status `json:"status"`

	updateTime time.Time
}

func NewRoom(name string) (*Room, error) {
	if name == "" {
		return nil, ErrEmptyRoomName
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	players := make([]*Player, 0, 5)
	return &Room{
		ID:         id,
		Name:       name,
		Players:    players,
		Status:     voting.StatusVoting,
		updateTime: time.Now(),
	}, nil
}

func (r *Room) AddPlayer(player *Player) error {
	r.Players = append(r.Players, player)
	return nil
}

func (r *Room) RemovePlayer(ID uuid.UUID) error {
	for index, player := range r.Players {
		if player.ID == ID {
			r.Players = append(r.Players[:index], r.Players[index+1:]...)
			return nil
		}
	}
	return ErrPlayerNotFound
}

func (r *Room) PlayerByID(ID uuid.UUID) (*Player, error) {
	for _, player := range r.Players {
		if player.ID == ID {
			return player, nil
		}
	}

	return nil, ErrPlayerNotFound
}

func (r *Room) RemoveAllPlayers() error {
	r.Players = make([]*Player, 0)
	return nil
}

func (r *Room) ClearVotes() error {
	for _, player := range r.Players {
		err := player.ClearVote()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Room) SetStatus(status voting.Status) error {
	r.Status = status
	return nil
}

func (r *Room) UpdateTime() time.Time {
	return r.updateTime
}
