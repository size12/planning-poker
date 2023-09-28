package handlers

import (
	"github.com/size12/planning-poker/internal/entity/voting"
	"github.com/size12/planning-poker/internal/service/room"
)

func (h *Handlers) CreateRoom(name string) (*room.Room, error) {
	r, err := room.New(name)
	if err != nil {
		return nil, err
	}

	r.SetButtonsPack(voting.PackModifiedFibo)

	h.pool.Store(r.ID, r)

	return r, nil
}
