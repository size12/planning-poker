package handlers

import (
	"github.com/google/uuid"
	"github.com/size12/planning-poker/internal/entity"
	"github.com/size12/planning-poker/internal/entity/voting"
	"github.com/size12/planning-poker/internal/service"
)

type Handlers struct {
	pool *service.RoomPool
}

func NewHandlers(pool *service.RoomPool) (*Handlers, error) {
	return &Handlers{pool: pool}, nil
}

func (h *Handlers) CreateRoom(name string) (*entity.Room, error) {
	room, err := entity.NewRoom(name)
	if err != nil {
		return nil, err
	}

	err = h.pool.AddRoom(room)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (h *Handlers) GetRoom(id uuid.UUID) (*entity.Room, error) {
	return h.pool.RoomByID(id)
}

func (h *Handlers) Vote(roomID, playerID uuid.UUID, vote *voting.Vote) error {
	room, err := h.pool.RoomByID(roomID)
	if err != nil {
		return err
	}

	room.UpdateTime()

	player, err := room.PlayerByID(playerID)
	if err != nil {
		return err
	}

	return player.Vote(vote)
}

func (h *Handlers) Reveal(roomID uuid.UUID) error {
	room, err := h.pool.RoomByID(roomID)
	if err != nil {
		return err
	}

	room.UpdateTime()

	return room.SetStatus(voting.StatusRevealed)
}

func (h *Handlers) StartVoting(roomID uuid.UUID) error {
	room, err := h.pool.RoomByID(roomID)
	if err != nil {
		return err
	}

	room.UpdateTime()

	return room.SetStatus(voting.StatusVoting)
}
