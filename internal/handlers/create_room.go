package handlers

import (
	"github.com/size12/planning-poker/internal/entity/voting"
	"github.com/size12/planning-poker/internal/service/room"
)

func (h *Handlers) CreateRoom(name string, buttons string) (*room.Room, error) {
	r, err := room.New(name)
	if err != nil {
		return nil, err
	}

	switch buttons {
	case "MODIFIED_FIBO":
		r.SetButtonsPack(voting.PackModifiedFibo)
	case "T_SHIRT":
		r.SetButtonsPack(voting.PackTShirt)
	case "SEQUENCE":
		r.SetButtonsPack(voting.PackSequence)
	default:
		r.SetButtonsPack(voting.PackModifiedFibo)
	}

	h.pool.Store(r.ID, r)

	return r, nil
}
