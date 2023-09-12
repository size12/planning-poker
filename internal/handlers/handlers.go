package handlers

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/size12/planning-poker/internal/service"
)

var ErrRoomNotFound = errors.New("room with such ID does not exists")

type Handlers struct {
	pool *sync.Map
}

func NewHandlers() (*Handlers, error) {
	pool := &sync.Map{}

	return &Handlers{pool: pool}, nil
}

func (h *Handlers) CreateRoom(name string) (*service.Room, error) {
	room, err := service.NewRoom(name)
	if err != nil {
		return nil, err
	}

	h.pool.Store(room.ID, room)

	return room, nil
}

func (h *Handlers) GetRoom(id uuid.UUID) (*service.Room, error) {
	value, ok := h.pool.Load(id)
	if !ok {
		return nil, ErrRoomNotFound
	}

	room, ok := value.(*service.Room)
	if !ok {
		return nil, fmt.Errorf("failed get *service.Room object from sync.Map")
	}

	return room, nil
}
