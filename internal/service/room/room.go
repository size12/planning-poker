package room

import (
	"sync"

	"github.com/google/uuid"
	"github.com/size12/planning-poker/internal/entity"
	"github.com/size12/planning-poker/internal/entity/voting"
)

type Room struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	Players []*entity.Player `json:"players"`
	Status  voting.Status    `json:"status"`
	History []entity.Score   `json:"history"`
	Buttons voting.VotePack  `json:"buttons"`
	Score   string           `json:"score"`

	adminID uuid.UUID

	*sync.RWMutex
}

func New(name string) (*Room, error) {
	if name == "" {
		return nil, ErrEmptyRoomName
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	hist := make([]entity.Score, 0, 5)

	players := make([]*entity.Player, 0, 5)
	return &Room{
		ID:      id,
		Name:    name,
		Players: players,
		Status:  voting.RoomStatusWaiting,
		RWMutex: &sync.RWMutex{},
		History: hist,
	}, nil
}

func (r *Room) Size() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.Players)
}
