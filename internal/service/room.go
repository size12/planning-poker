package service

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/size12/planning-poker/internal/entity"
	"github.com/size12/planning-poker/internal/entity/voting"
)

var ErrPlayerNotFound = errors.New("player with such ID does not exists")
var ErrEmptyRoomName = errors.New("room can't have empty name")
var ErrPlayerAlreadyExists = errors.New("room already have this player")

type Room struct {
	ID      uuid.UUID        `json:"id"`
	Name    string           `json:"name"`
	Players []*entity.Player `json:"players"`
	Status  voting.Status    `json:"status"`
	AdminID uuid.UUID        `json:"-"`

	*sync.RWMutex
}

type RoomInfo struct {
	ID      uuid.UUID `json:"id"`
	AdminID uuid.UUID `json:"admin_id"`
}

func NewRoom(name string) (*Room, error) {
	if name == "" {
		return nil, ErrEmptyRoomName
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	players := make([]*entity.Player, 0, 5)
	return &Room{
		ID:      id,
		Name:    name,
		Players: players,
		Status:  voting.StatusVoting,
		RWMutex: &sync.RWMutex{},
	}, nil
}

func (r *Room) AddPlayer(player *entity.Player) error {
	r.Lock()
	defer r.Unlock()

	for _, p := range r.Players {
		if p.ID == player.ID {
			return ErrPlayerAlreadyExists
		}
	}

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

func (r *Room) ClearVotes() error {
	r.Lock()
	defer r.Unlock()

	for _, player := range r.Players {
		err := player.ClearVote()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Room) SetStatus(status voting.Status) error {
	r.Lock()
	defer r.Unlock()

	r.Status = status
	return nil
}
