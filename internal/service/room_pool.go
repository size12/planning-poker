package service

//
//import (
//	"errors"
//	"log"
//	"sync"
//
//	"github.com/google/uuid"
//	"github.com/size12/planning-poker/internal/entity"
//)
//

//
//type RoomPool struct {
//	rooms *sync.Map
//}
//
//func NewRoomPool() (*RoomPool, error) {
//	rooms := &sync.Map{}
//	return &RoomPool{rooms: rooms}, nil
//}
//
//func (p *RoomPool) RoomByID(ID uuid.UUID) (*entity.Room, error) {
//	loaded, ok := p.rooms.Load(ID)
//	if !ok {
//		return nil, ErrRoomNotFound
//	}
//
//	room, ok := loaded.(*entity.Room)
//	if !ok {
//		log.Fatalln("There is not room object in map.")
//	}
//
//	return room, nil
//}
//
//func (p *RoomPool) AddRoom(room *entity.Room) error {
//	p.rooms.Store(room.ID, room)
//	return nil
//}
//
//func (p *RoomPool) RemoveRoom(id uuid.UUID) error {
//	p.rooms.Delete(id)
//	return nil
//}
