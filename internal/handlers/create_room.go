package handlers

import (
	"github.com/size12/planning-poker/internal/service/room"
)

func (h *Handlers) CreateRoom(name string) (*room.Room, error) {
	r, err := room.New(name)
	if err != nil {
		return nil, err
	}

	h.pool.Store(r.ID, r)

	return r, nil
}
