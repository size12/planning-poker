package websocket

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/olahol/melody"
	"github.com/size12/planning-poker/internal/service/room"
)

func WsDisconnected(r *room.Room, m *melody.Melody) func(*melody.Session) {
	return func(session *melody.Session) {
		object, ok := session.Get("player_id")
		if !ok {
			log.Printf("Someone without user_id disconnected from room %s\n", r.ID)
			return
		}

		ID, ok := object.(uuid.UUID)
		if !ok {
			log.Printf("Someone without valid user_id disconnected from room %s\n", r.ID)
			return
		}

		err := r.RemovePlayer(ID)
		if err != nil {
			log.Printf("Failed remove player %s from room %s: %v\n", ID, r.ID, err)
			return
		}

		log.Printf("User %s disconnected from room %s\n", ID, r.ID)

		r.RLock()
		defer r.RUnlock()

		bytes, err := json.Marshal(r)
		if err != nil {
			log.Printf("Failed marshal response: %v\n", err)
			return
		}

		err = m.BroadcastOthers(bytes, session)
		if err != nil {
			log.Printf("Failed broadcast room %s: %v\n", r.ID, err)
			return
		}
	}
}
