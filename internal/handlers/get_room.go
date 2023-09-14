package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/size12/planning-poker/internal/service/room"
)

func (h *Handlers) RoomByID(id uuid.UUID) (*room.Room, error) {
	value, ok := h.pool.Load(id)
	if !ok {
		return nil, ErrRoomNotFound
	}

	r, ok := value.(*room.Room)
	if !ok {
		return nil, fmt.Errorf("failed get *service.CreateRoom object from sync.Map")
	}

	return r, nil
}
