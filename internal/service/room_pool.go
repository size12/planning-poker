package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/size12/planning-poker/internal/entity"
)

var ErrRoomNotFound = errors.New("room with such ID does not exists")

type RoomPool struct {
	rooms map[uuid.UUID]*entity.Room
}

func NewRoomPool() (*RoomPool, error) {
	rooms := make(map[uuid.UUID]*entity.Room)
	return &RoomPool{rooms: rooms}, nil
}

func (p *RoomPool) RoomByID(ID uuid.UUID) (*entity.Room, error) {
	room, ok := p.rooms[ID]
	if !ok {
		return nil, ErrRoomNotFound
	}

	return room, nil
}

func (p *RoomPool) AddRoom(room *entity.Room) error {
	p.rooms[room.ID] = room
	return nil
}

func (p *RoomPool) RemoveRoom(id uuid.UUID) error {
	delete(p.rooms, id)
	return nil
}

func (p *RoomPool) DeleteInactive(ctx context.Context, timeout time.Duration) {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			for _, room := range p.rooms {
				if room.UpdateTime().Sub(now) < timeout {
					err := p.RemoveRoom(room.ID)
					if err != nil {
						log.Printf("Failed delete inactive room with ID %s: %v\n", room.ID, err)
					}
				}
			}
		case <-ctx.Done():
			log.Println("Shutdown delete inactive rooms service.")
			return
		}
	}
}
